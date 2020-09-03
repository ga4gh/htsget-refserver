package htsserver

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var httpRequestSingleTC = []struct {
	endpoint string
	query    [][]string
	expCode  int
	expBody  string
}{
	/* SERVICE-INFO TEST CASES */
	{
		"/reads/service-info",
		nil,
		200,
		"{\"id\":\"htsgetref.reads\",\"name\":\"GA4GH htsget reference server reads endpoint\",\"type\":{\"group\":\"org.ga4gh\",\"artifact\":\"htsget\",\"version\":\"1.2.0\"},\"description\":\"Stream alignment files (BAM/CRAM) according to GA4GH htsget protocol\",\"organization\":{\"name\":\"Global Alliance for Genomics and Health\",\"url\":\"https://ga4gh.org\"},\"contactUrl\":\"mailto:jeremy.adams@ga4gh.org\",\"documentationUrl\":\"https://ga4gh.org\",\"createdAt\":\"2020-09-01T12:00:00Z\",\"updatedAt\":\"2020-09-01T12:00:00Z\",\"environment\":\"test\",\"version\":\"1.3.0\",\"htsget\":{\"datatype\":\"reads\",\"formats\":[\"BAM\"],\"fieldsParameterEffective\":true,\"tagsParametersEffective\":true}}\n",
	},
	{
		"/variants/service-info",
		nil,
		200,
		"{\"id\":\"htsgetref.variants\",\"name\":\"GA4GH htsget reference server variants endpoint\",\"type\":{\"group\":\"org.ga4gh\",\"artifact\":\"htsget\",\"version\":\"1.2.0\"},\"description\":\"Stream variant files (VCF/BCF) according to GA4GH htsget protocol\",\"organization\":{\"name\":\"Global Alliance for Genomics and Health\",\"url\":\"https://ga4gh.org\"},\"contactUrl\":\"mailto:jeremy.adams@ga4gh.org\",\"documentationUrl\":\"https://ga4gh.org\",\"createdAt\":\"2020-09-01T12:00:00Z\",\"updatedAt\":\"2020-09-01T12:00:00Z\",\"environment\":\"test\",\"version\":\"1.3.0\",\"htsget\":{\"datatype\":\"variants\",\"formats\":[\"VCF\"],\"fieldsParameterEffective\":true,\"tagsParametersEffective\":true}}\n",
	},
	/* READS TICKET CASES */
	{
		"/reads/NonExistentId",
		nil,
		404,
		"{\"htsget\":{\"error\":\"NotFound\",\"message\":\"The requested resource could not be associated with a registered data source\"}}\n",
	},
	{
		"/reads/tabulamuris.object00001",
		nil,
		404,
		"{\"htsget\":{\"error\":\"NotFound\",\"message\":\"The requested resource was not found\"}}\n",
	},
	{
		"/reads/tabulamuris.A1-B000168-3_57_F-1-1_R2",
		nil,
		200,
		"{\"htsget\":{\"format\":\"BAM\",\"urls\":[{\"url\":\"https://s3.amazonaws.com/czbiohub-tabula-muris/facs_bam_files/A1-B000168-3_57_F-1-1_R2.mus.Aligned.out.sorted.bam\",\"headers\":{\"Range\":\"bytes=0-41157\"}}]}}\n",
	},
	{
		"/reads/tabulamuris.A1-B000168-3_57_F-1-1_R2",
		[][]string{
			[]string{"class", "header"},
		},
		200,
		"{\"htsget\":{\"format\":\"BAM\",\"urls\":[{\"url\":\"http://localhost:3000/reads/data/tabulamuris.A1-B000168-3_57_F-1-1_R2?class=header\",\"headers\":{\"HtsgetBlockId\":\"1\",\"HtsgetNumBlocks\":\"1\",\"HtsgetBlockClass\":\"header\"},\"class\":\"header\"}]}}\n",
	},
	{
		"/reads/tabulamuris.A1-B000168-3_57_F-1-1_R2",
		[][]string{
			[]string{"format", "BAM"},
			[]string{"referenceName", "chr1"},
			[]string{"start", "20000000"},
			[]string{"end", "30000000"},
			[]string{"fields", "SEQ,QUAL"},
			[]string{"tags", "HI,NM"},
		},
		200,
		"{\"htsget\":{\"format\":\"BAM\",\"urls\":[{\"url\":\"http://localhost:3000/reads/data/tabulamuris.A1-B000168-3_57_F-1-1_R2?end=30000000\\u0026fields=SEQ%2CQUAL\\u0026referenceName=chr1\\u0026start=20000000\\u0026tags=HI%2CNM\",\"headers\":{\"HtsgetBlockId\":\"1\",\"HtsgetNumBlocks\":\"2\",\"HtsgetBlockClass\":\"header\"},\"class\":\"header\"},{\"url\":\"http://localhost:3000/reads/data/tabulamuris.A1-B000168-3_57_F-1-1_R2?end=30000000\\u0026fields=SEQ%2CQUAL\\u0026referenceName=chr1\\u0026start=20000000\\u0026tags=HI%2CNM\",\"headers\":{\"HtsgetBlockId\":\"2\",\"HtsgetNumBlocks\":\"2\"}}]}}\n",
	},
}

func TestHTTPRequestSingle(t *testing.T) {
	router, _ := SetRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	for _, tc := range httpRequestSingleTC {
		// construct url
		url := server.URL + tc.endpoint

		// assign query params to the url
		if tc.query != nil {
			queryParams := []string{}
			for _, p := range tc.query {
				queryParams = append(queryParams, p[0]+"="+p[1])
			}
			url += "?" + strings.Join(queryParams, "&")
		}

		// make GET request to server
		resp, _ := http.Get(url)

		// read response body into byte array
		responseBodyBytes := make([]byte, 4096)
		nBytes, _ := resp.Body.Read(responseBodyBytes)
		responseBody := string(responseBodyBytes[0:nBytes])

		// assert status code, response body
		assert.Equal(t, tc.expCode, resp.StatusCode)
		assert.Equal(t, tc.expBody, responseBody)
	}
}
