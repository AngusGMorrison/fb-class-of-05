package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

type mockHTTPLogger struct {
	wasCalled bool
}

func (m *mockHTTPLogger) Log(ww *WrappedWriter, r *http.Request, d time.Duration) {
	m.wasCalled = true
}

func TestLogging(t *testing.T) {
	var logger mockHTTPLogger
	mux := mux.NewRouter()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		return
	})
	mux.Use(Logging(&logger))

	req, _ := http.NewRequest("GET", "/", nil)
	mux.ServeHTTP(httptest.NewRecorder(), req)
	if logger.wasCalled != true {
		t.Errorf("logger was not called by Logging middleware")
	}
}
