package assumerole

import (
	"github.com/ga4gh/htsget-refserver/internal/awsutils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var mockHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("bar"))
})

// go test -run TestNew ./internal/assumerole/ -v -count 1
func TestNew(t *testing.T) {
	ar := New(Options{Debug: true})
	assert.True(t, ar != nil)
}

// go test -run TestHandler ./internal/assumerole/ -v -count 1
func TestHandler(t *testing.T) {

	if os.Getenv(awsutils.AwsProfile) != awsutils.TestAwsProfileForIT {
		t.Skipf("[Skip] Required to setup and `export AWS_PROFILE=%s` in integration testing CI environment", awsutils.TestAwsProfileForIT)
	}

	h := Handler(Options{Debug: true})
	req, _ := http.NewRequest("GET", "http://localhost/reads/foo", nil)
	res := httptest.NewRecorder()
	h(mockHandler).ServeHTTP(res, req)

	val, exist := os.LookupEnv(awsutils.AwsAccessKeyId)
	assert.True(t, exist)
	assert.True(t, len(val) > 0)
}
