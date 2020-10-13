// Package htsrequest provides operations for parsing htsget-related
// parameters from the HTTP request
//
// Module parsing_test tests module parsing
package htsrequest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

// parsePathParamTC test cases for parsePathParam
var parsePathParamTC = []struct {
	id       string
	expValue string
	expFound bool
}{
	{"readsId12345", "readsId12345", true},
	{"object0001", "object0001", true},
	{"", "", false},
}

// parseQueryParamTC test cases for parseQueryParam
var parseQueryParamTC = []struct {
	url, key string
	expValue string
	expFound bool
	expError bool
}{
	{
		"/reads/01?format=BAM",
		"format",
		"BAM",
		true,
		false,
	},
	{
		"/reads/01?class=header",
		"format",
		"",
		false,
		false,
	},
	{
		"/reads/01?format=BAM&format=CRAM",
		"format",
		"",
		false,
		true,
	},
}

// parseHeaderParamTC test cases for parseHeaderParam
var parseHeaderParamTC = []struct {
	headers  [][]string
	key      string
	expValue string
	expFound bool
}{
	{
		[][]string{
			[]string{"Range", "bytes=1000-2000"},
		},
		"Range",
		"bytes=1000-2000",
		true,
	},
	{
		[][]string{
			[]string{"HtsgetTotalBlocks", "100"},
		},
		"Range",
		"",
		false,
	},
}

// parseReqBodyParamTC test cases for parseReqBodyParam
var parseReqBodyParamTC = []struct {
	rawBody, key       string
	expError, expFound bool
	expValue           string
}{}

// TestParsePathParam tests parsePathParam function
func TestParsePathParam(t *testing.T) {
	for _, tc := range parsePathParamTC {
		// mock router to parse path param and assert its existence within
		// the handler function
		router := chi.NewRouter()
		router.Get("/test/{id}", func(writer http.ResponseWriter, request *http.Request) {
			value, found := parsePathParam(request, "id")
			assert.Equal(t, tc.expFound, found)
			if tc.expFound {
				assert.Equal(t, tc.expValue, value)
			}
		})

		request := httptest.NewRequest("GET", "/test/"+tc.id, nil)
		writer := httptest.NewRecorder()
		router.ServeHTTP(writer, request)
	}
}

// TestParseQueryParam tests parseQueryParam function
func TestParseQueryParam(t *testing.T) {
	for _, tc := range parseQueryParamTC {
		request := httptest.NewRequest("GET", tc.url, nil)
		value, found, err := parseQueryParam(request.URL.Query(), tc.key)
		if tc.expError {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tc.expFound, found)
			assert.Equal(t, tc.expValue, value)
		}
	}
}

// TestParseHeaderParam tests parseHeaderParam function
func TestParseHeaderParam(t *testing.T) {
	for _, tc := range parseHeaderParamTC {
		request := httptest.NewRequest("GET", "/test/object01", nil)
		for _, header := range tc.headers {
			request.Header.Add(header[0], header[1])
		}
		value, found := parseHeaderParam(request, tc.key)
		assert.Equal(t, tc.expFound, found)
		assert.Equal(t, tc.expValue, value)
	}
}

func TestParseReqBodyParam(t *testing.T) {
	// for _, tc := range parseReqBodyParamTC {
	//
	// }
}
