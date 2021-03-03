package assumerole

import (
	"github.com/ga4gh/htsget-refserver/internal/awsutils"
	"log"
	"net/http"
	"os"
)

type Options struct {
	Debug bool
}

type Logger interface {
	Printf(string, ...interface{})
}

type AssumeRole struct {
	Log Logger
}

func New(options Options) *AssumeRole {
	ar := &AssumeRole{}
	if options.Debug && ar.Log == nil {
		ar.Log = log.New(os.Stdout, "[assumerole] ", log.LstdFlags)
	}
	return ar
}

func Handler(options Options) func(next http.Handler) http.Handler {
	ar := New(options)
	return ar.Handler
}

func (ar *AssumeRole) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		awsutils.UnsetCredentials()
		cred, err := awsutils.GetCredentials()
		if err != nil {
			ar.logf("error getting credentials with assume role")
			ar.logf(err.Error())
		} else {
			ar.logf("using credentials source " + cred.Source)
			awsutils.SetCredentials(*cred)
		}
		next.ServeHTTP(w, r)
	})
}

func (ar *AssumeRole) logf(format string, a ...interface{}) {
	if ar.Log != nil {
		ar.Log.Printf(format, a...)
	}
}
