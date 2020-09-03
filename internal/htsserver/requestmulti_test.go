package htsserver

import (
	"bufio"
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

var httpRequestMultiTC = []struct {
	endpoint    string
	queryParams [][]string
	expFilename string
}{
	/* **************************************************
	 * TEST CASES
	 * ************************************************** */

	/* **************************************************
	 * READS
	 * ************************************************** */

	// READS, NOTHING SPECIFIED

	{
		"/reads/tabulamuris.A1-B000168-3_57_F-1-1_R2",
		nil,
		"reads-tc-00.bam",
	},
	// READS, SPECIFY REFERENCE NAME, START, END, FIELDS

	{
		"/reads/tabulamuris.A1-B000168-3_57_F-1-1_R2",
		[][]string{
			[]string{"referenceName", "chr1"},
			[]string{"start", "20000000"},
			[]string{"end", "30000000"},
			[]string{"fields", "SEQ,QUAL"},
		},
		"reads-tc-01.bam",
	},

	/* **************************************************
	 * VARIANTS
	 * ************************************************** */
}

func getHtsgetTicket(server *httptest.Server, endpoint string, queryParams [][]string) *htsticket.Ticket {
	url := server.URL + endpoint
	url = assignQueryParams(url, queryParams)
	resp, _ := http.Get(url)

	responseBodyBytesAll := make([]byte, 4096)
	nBytes, _ := resp.Body.Read(responseBodyBytesAll)
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
	headerKeys := []string{"HtsgetBlockId", "HtsgetNumBlocks", "Range", "HtsgetBlockClass", "HtsgetFilePath"}
	headerVals := []string{h.BlockID, h.NumBlocks, h.Range, h.Class, h.FilePath}
	for a := range headerKeys {
		if headerVals[a] != "" {
			request.Header.Set(headerKeys[a], headerVals[a])
		}
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
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
		ticket := getHtsgetTicket(server, tc.endpoint, tc.queryParams)

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
