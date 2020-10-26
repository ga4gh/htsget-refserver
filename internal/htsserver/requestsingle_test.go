package htsserver

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ga4gh/htsget-refserver/internal/htsconfig"
	"github.com/stretchr/testify/assert"
)

var httpRequestSingleTC = []struct {
	method      string
	endpoint    string
	query       [][]string
	headers     [][]string
	requestBody string
	expCode     int
	expBody     string
}{
	/* GET SERVICE-INFO TEST CASES */
	{
		"GET",
		"/reads/service-info",
		nil,
		nil,
		"",
		200,
		"{\"id\":\"htsgetref.reads\",\"name\":\"GA4GH htsget reference server reads endpoint\",\"type\":{\"group\":\"org.ga4gh\",\"artifact\":\"htsget\",\"version\":\"1.2.0\"},\"description\":\"Stream alignment files (BAM/CRAM) according to GA4GH htsget protocol\",\"organization\":{\"name\":\"Global Alliance for Genomics and Health\",\"url\":\"https://ga4gh.org\"},\"contactUrl\":\"mailto:jeremy.adams@ga4gh.org\",\"documentationUrl\":\"https://ga4gh.org\",\"createdAt\":\"2020-09-01T12:00:00Z\",\"updatedAt\":\"2020-09-01T12:00:00Z\",\"environment\":\"test\",\"version\":\"1.4.1\",\"htsget\":{\"datatype\":\"reads\",\"formats\":[\"BAM\"],\"fieldsParameterEffective\":true,\"tagsParametersEffective\":true}}\n",
	},
	{
		"GET",
		"/variants/service-info",
		nil,
		nil,
		"",
		200,
		"{\"id\":\"htsgetref.variants\",\"name\":\"GA4GH htsget reference server variants endpoint\",\"type\":{\"group\":\"org.ga4gh\",\"artifact\":\"htsget\",\"version\":\"1.2.0\"},\"description\":\"Stream variant files (VCF/BCF) according to GA4GH htsget protocol\",\"organization\":{\"name\":\"Global Alliance for Genomics and Health\",\"url\":\"https://ga4gh.org\"},\"contactUrl\":\"mailto:jeremy.adams@ga4gh.org\",\"documentationUrl\":\"https://ga4gh.org\",\"createdAt\":\"2020-09-01T12:00:00Z\",\"updatedAt\":\"2020-09-01T12:00:00Z\",\"environment\":\"test\",\"version\":\"1.4.1\",\"htsget\":{\"datatype\":\"variants\",\"formats\":[\"VCF\"],\"fieldsParameterEffective\":false,\"tagsParametersEffective\":false}}\n",
	},
	/* GET READS TICKET CASES */
	{
		"GET",
		"/reads/NonExistentId",
		nil,
		nil,
		"",
		404,
		"{\"htsget\":{\"error\":\"NotFound\",\"message\":\"The requested resource could not be associated with a registered data source\"}}\n",
	},
	{
		"GET",
		"/reads/tabulamuris.object00001",
		nil,
		nil,
		"",
		404,
		"{\"htsget\":{\"error\":\"NotFound\",\"message\":\"The requested resource was not found\"}}\n",
	},

	{
		"GET",
		"/reads/tabulamuris.A1-B000168-3_57_F-1-1_R2",
		nil,
		nil,
		"",
		200,
		"{\"htsget\":{\"format\":\"BAM\",\"urls\":[{\"url\":\"http://localhost:3000/file-bytes\",\"headers\":{\"HtsgetFilePath\":\"../../data/test/sources/tabulamuris/A1-B000168-3_57_F-1-1_R2.mus.Aligned.out.sorted.bam\",\"Range\":\"bytes=0-41157\"}}]}}\n",
	},

	{
		"GET",
		"/reads/tabulamuris.A1-B000168-3_57_F-1-1_R2",
		[][]string{
			[]string{"class", "header"},
		},
		nil,
		"",
		200,
		"{\"htsget\":{\"format\":\"BAM\",\"urls\":[{\"url\":\"http://localhost:3000/reads/data/tabulamuris.A1-B000168-3_57_F-1-1_R2?class=header\",\"headers\":{\"HtsgetBlockClass\":\"header\",\"HtsgetCurrentBlock\":\"0\",\"HtsgetTotalBlocks\":\"1\"},\"class\":\"header\"}]}}\n",
	},

	{
		"GET",
		"/reads/tabulamuris.A1-B000168-3_57_F-1-1_R2",
		[][]string{
			[]string{"format", "BAM"},
			[]string{"referenceName", "chr1"},
			[]string{"start", "20000000"},
			[]string{"end", "30000000"},
			[]string{"fields", "SEQ,QUAL"},
			[]string{"tags", "HI,NM"},
		},
		nil,
		"",
		200,
		"{\"htsget\":{\"format\":\"BAM\",\"urls\":[{\"url\":\"http://localhost:3000/reads/data/tabulamuris.A1-B000168-3_57_F-1-1_R2?fields=SEQ%2CQUAL\\u0026tags=HI%2CNM\",\"headers\":{\"HtsgetBlockClass\":\"header\",\"HtsgetCurrentBlock\":\"0\",\"HtsgetTotalBlocks\":\"2\"},\"class\":\"header\"},{\"url\":\"http://localhost:3000/reads/data/tabulamuris.A1-B000168-3_57_F-1-1_R2?end=30000000\\u0026fields=SEQ%2CQUAL\\u0026referenceName=chr1\\u0026start=20000000\\u0026tags=HI%2CNM\",\"headers\":{\"HtsgetCurrentBlock\":\"1\",\"HtsgetTotalBlocks\":\"2\"},\"class\":\"body\"}]}}\n",
	},

	{
		"GET",
		"/reads/tabulamuris.A1-B000168-3_57_F-1-1_R2",
		[][]string{
			[]string{"format", "BAM"},
			[]string{"fields", "QNAME,FLAG,RNAME"},
			[]string{"tags", "HI,NM"},
		},
		nil,
		"",
		200,
		"{\"htsget\":{\"format\":\"BAM\",\"urls\":[{\"url\":\"http://localhost:3000/reads/data/tabulamuris.A1-B000168-3_57_F-1-1_R2?fields=QNAME%2CFLAG%2CRNAME\\u0026tags=HI%2CNM\",\"headers\":{\"HtsgetBlockClass\":\"header\",\"HtsgetCurrentBlock\":\"0\",\"HtsgetTotalBlocks\":\"2\"},\"class\":\"header\"},{\"url\":\"http://localhost:3000/reads/data/tabulamuris.A1-B000168-3_57_F-1-1_R2?fields=QNAME%2CFLAG%2CRNAME\\u0026tags=HI%2CNM\",\"headers\":{\"HtsgetCurrentBlock\":\"1\",\"HtsgetTotalBlocks\":\"2\"},\"class\":\"body\"}]}}\n",
	},

	/* GET VARIANTS TICKET CASES */

	{
		"GET",
		"/variants/HG002_GIAB",
		nil,
		nil,
		"",
		200,
		"{\"htsget\":{\"format\":\"VCF\",\"urls\":[{\"url\":\"http://localhost:3000/file-bytes\",\"headers\":{\"HtsgetFilePath\":\"../../data/test/sources/giab/HG002_GIAB.filtered.vcf.gz\",\"Range\":\"bytes=0-185237\"}}]}}\n",
	},
}

func TestHTTPRequestSingle(t *testing.T) {

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

	router, _ := SetRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	for _, tc := range httpRequestSingleTC {
		// construct url
		url := server.URL + tc.endpoint

		// query params
		if tc.query != nil {
			queryParams := []string{}
			for _, p := range tc.query {
				queryParams = append(queryParams, p[0]+"="+p[1])
			}
			url += "?" + strings.Join(queryParams, "&")
		}

		// requestBody
		var requestBody io.Reader = nil

		// instantiate test http request
		request := httptest.NewRequest(tc.method, url, requestBody)

		// headers
		if tc.headers != nil {
			for _, header := range tc.headers {
				request.Header.Add(header[0], header[1])
			}
		}

		// write response to response recorder and convert response body to string
		writer := httptest.NewRecorder()
		router.ServeHTTP(writer, request)
		responseBodyBytes, _ := ioutil.ReadAll(writer.Body)
		responseBody := string(responseBodyBytes)

		// assert status code, response body
		assert.Equal(t, tc.expCode, writer.Code)
		assert.Equal(t, tc.expBody, responseBody)
	}

	// set the configuration back to default
	htsconfig.SetConfigFile(htsconfig.DefaultConfiguration)
	htsconfig.SetConfig(htsconfig.DefaultConfiguration)
}
