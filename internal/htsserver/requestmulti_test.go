package htsserver

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ga4gh/htsget-refserver/internal/htsconfig"

	"github.com/ga4gh/htsget-refserver/internal/htsticket"
	"github.com/stretchr/testify/assert"
)

type E2ETestCase struct {
	method      string
	endpoint    string
	queryParams [][]string
	requestBody string
	expFilename string
}

var httpRequestMultiTC = []E2ETestCase{
	/* **************************************************
	 * TEST CASES
	 * ************************************************** */

	/* **************************************************
	 * READS
	 * ************************************************** */

	// GET READS, NOTHING SPECIFIED

	{
		"GET",
		"/reads/tabulamuris.A1-B000168-3_57_F-1-1_R2",
		nil,
		"",
		"reads-tc-00.bam",
	},

	// GET READS, SPECIFY REFERENCE NAME, START, END, FIELDS

	{
		"GET",
		"/reads/tabulamuris.A1-B000168-3_57_F-1-1_R2",
		[][]string{
			[]string{"referenceName", "chr1"},
			[]string{"start", "20000000"},
			[]string{"end", "30000000"},
			[]string{"fields", "SEQ,QUAL"},
		},
		"",
		"reads-tc-01.bam",
	},

	// GET READS, SPECIFY REGION
	{
		"GET",
		"/reads/tabulamuris.A1-B000168-3_57_F-1-1_R2",
		[][]string{
			[]string{"referenceName", "chr1"},
			[]string{"start", "20000000"},
			[]string{"end", "30000000"},
		},
		"",
		"reads-tc-02.bam",
	},

	// GET READS, SPECIFY FIELDS, TAGS
	{
		"GET",
		"/reads/tabulamuris.A1-B000168-3_57_F-1-1_R2",
		[][]string{
			[]string{"fields", "QNAME,FLAG,SEQ,QUAL"},
			[]string{"tags", "HI,NM"},
		},
		"",
		"reads-tc-03.bam",
	},

	// GET READS, SPECIFY TAGS=""
	{
		"GET",
		"/reads/tabulamuris.A1-B000168-3_57_F-1-1_R2",
		[][]string{
			[]string{"tags", ""},
		},
		"",
		"reads-tc-04.bam",
	},

	// GET READS, SPECIFY NOTAGS
	{
		"GET",
		"/reads/tabulamuris.A1-B000168-3_57_F-1-1_R2",
		[][]string{
			[]string{"fields", "QNAME,FLAG,SEQ,QUAL"},
			[]string{"notags", "HI,NM"},
		},
		"",
		"reads-tc-05.bam",
	},

	// POST READS, SPECIFY REGIONS

	{
		"POST",
		"/reads/tabulamuris.A1-B000168-3_57_F-1-1_R2",
		nil,
		"{\"regions\":[{\"referenceName\":\"chr10\"},{\"referenceName\":\"chr13\"},{\"referenceName\":\"chr16\"}]}",
		"reads-tc-06.bam",
	},

	// POST READS, SPECIFY REGIONS, FIELDS
	{
		"POST",
		"/reads/tabulamuris.A1-B000168-3_57_F-1-1_R2",
		nil,
		"{\"fields\":[\"QNAME\",\"FLAG\",\"RNAME\",\"POS\"],\"regions\":[{\"referenceName\":\"chr7\"},{\"referenceName\":\"chr11\"}]}",
		"reads-tc-07.bam",
	},

	// POST READS, SPECIFY REGIONS, FIELDS, TAGS
	{
		"POST",
		"/reads/tabulamuris.A1-B000168-3_57_F-1-1_R2",
		nil,
		"{\"fields\":[\"RNAME\",\"POS\"],\"tags\":[\"MD\"],\"regions\":[{\"referenceName\":\"chr8\"},{\"referenceName\":\"chr12\"}]}",
		"reads-tc-08.bam",
	},

	// POST READS, SPECIFY REGIONS, FIELDS, NOTAGS
	{
		"POST",
		"/reads/tabulamuris.A1-B000168-3_57_F-1-1_R2",
		nil,
		"{\"fields\":[\"QNAME\",\"RNAME\",\"POS\",\"SEQ\",\"QUAL\"],\"notags\":[\"MD\"],\"regions\":[{\"referenceName\":\"chr5\"}]}",
		"reads-tc-09.bam",
	},

	/* **************************************************
	 * VARIANTS
	 * ************************************************** */

	// GET VARIANTS, NOTHING SPECIFIED

	{
		"GET",
		"/variants/HG002_GIAB",
		nil,
		"",
		"variants-tc-00.vcf.gz",
	},

	// GET VARIANTS, SPECIFY REFERENCE NAME

	{
		"GET",
		"/variants/HG002_GIAB",
		[][]string{
			[]string{"referenceName", "22"},
		},
		"",
		"variants-tc-01.vcf",
	},
}

func getHtsgetTicket(server *httptest.Server, tc E2ETestCase) *htsticket.Ticket {
	url := server.URL + tc.endpoint
	url = assignQueryParams(url, tc.queryParams)
	var response *http.Response
	if tc.method == "GET" {
		response, _ = http.Get(url)
	} else if tc.method == "POST" {
		response, _ = http.Post(url, "application/json", bytes.NewReader([]byte(tc.requestBody)))
	}

	responseBodyBytesAll := make([]byte, 4096)
	nBytes, _ := response.Body.Read(responseBodyBytesAll)
	responseBodyBytes := responseBodyBytesAll[0:nBytes]
	ticket := new(htsticket.Ticket)
	json.Unmarshal(responseBodyBytes, ticket)

	return ticket
}

func assignQueryParams(url string, queryParams [][]string) string {
	// assign query params to the url
	if queryParams != nil {
		pList := []string{}
		for _, p := range queryParams {
			pList = append(pList, p[0]+"="+p[1])
		}
		url += "?" + strings.Join(pList, "&")
	}
	return url
}

func downloadFilepart(server *httptest.Server, i int, ticketURL *htsticket.URL, writer *bufio.Writer) error {
	request, _ := http.NewRequest("GET", ticketURL.URL, nil)

	h := ticketURL.Headers
	headerKeys := []string{"HtsgetCurrentBlock", "HtsgetTotalBlocks", "Range", "HtsgetBlockClass", "HtsgetFilePath"}
	headerVals := []string{h.CurrentBlock, h.TotalBlocks, h.Range, h.BlockClass, h.FilePath}
	for a := range headerKeys {
		if headerVals[a] != "" {
			request.Header.Set(headerKeys[a], headerVals[a])
		}
	}

	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		return err
	}

	// defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	writer.Write(body)
	writer.Flush()
	return nil
}

func readBytesFromFile(fp string, size int) []byte {
	file, _ := os.Open(fp)
	reader := bufio.NewReader(file)
	bytesAll := make([]byte, size)
	nBytes, _ := reader.Read(bytesAll)
	bytes := bytesAll[0:nBytes]
	return bytes
}

func calculateMD5(fp string) string {
	bytes := readBytesFromFile(fp, 2097152)
	return fmt.Sprintf("%x", md5.Sum(bytes))
}

func TestHTTPRequestMulti(t *testing.T) {

	// configure dir in which temp and test comparator files are relative to
	wd, _ := os.Getwd()
	parentDir := filepath.Dir(filepath.Dir(wd))

	// set the configuration for E2E tests
	configFilePath := filepath.Join(parentDir, "data", "config", "integration-tests.config.json")
	configFile, _ := os.Open(configFilePath)
	configJSONBytes, _ := ioutil.ReadAll(configFile)
	newConfig := new(htsconfig.Configuration)
	json.Unmarshal(configJSONBytes, newConfig)
	htsconfig.SetConfigFile(newConfig)
	htsconfig.LoadConfig()

	// setup test server on port 3000
	router, _ := SetRouter()
	listener, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewUnstartedServer(router)
	server.Listener.Close()
	server.Listener = listener
	server.Start()
	defer server.Close()

	// run the htsget request loop for each test case
	for _, tc := range httpRequestMultiTC {
		// create the temp outputfile that htsget data response blocks will be
		// written to

		outputFilepath := htsconfig.GetTempfilePath("testoutput")
		outputFile, err := htsconfig.CreateTempfile("testoutput")
		if err != nil {
			t.Fatal(err)
		}

		// get the htsget ticket from an initial request
		writer := bufio.NewWriter(outputFile)
		ticket := getHtsgetTicket(server, tc)

		// for each url in the ticket's list of urls, download the filepart
		for i, ticketURL := range ticket.HTSget.URLS {
			err = downloadFilepart(server, i, ticketURL, writer)
			if err != nil {
				t.Fatal(err)
			}
		}
		outputFile.Close()

		// compare md5sum of the fileparts concatenated together against the
		// expected file in the test data directory
		expectedFilePath := filepath.Join(parentDir, "data", "test", "expected", tc.expFilename)
		expectedMD5 := calculateMD5(expectedFilePath)
		actualMD5 := calculateMD5(outputFilepath)
		htsconfig.RemoveTempfile(outputFile)
		assert.Equal(t, expectedMD5, actualMD5)
	}

	// set the configuration back to default
	htsconfig.SetConfigFile(htsconfig.DefaultConfiguration)
	htsconfig.SetConfig(htsconfig.DefaultConfiguration)
}
