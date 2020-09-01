package middleware

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type mockStackLogger struct {
	bytes.Buffer
}

func (m *mockStackLogger) LogStructuredStack(tracer StackTracer) {
	for _, f := range tracer.StackTrace() {
		m.WriteString(fmt.Sprintf("%+s:%d\n", f, f))
	}
}

func (m *mockStackLogger) LogUnstructuredStack(err error) {
	buf := make([]byte, 1000)
	runtime.Stack(buf, false)
	m.Write(buf)
}

func TestRecovery(t *testing.T) {
	testCases := []struct {
		desc                string
		path                string
		hf                  http.HandlerFunc
		wantOutputToContain string
	}{
		{
			"recovered object implements stackTracer",
			"/",
			func(w http.ResponseWriter, r *http.Request) {
				panic(errors.New("test error"))
			},
			"TestRecovery.func1", // the anonymous http.HandlerFunc where the trace is created
		},
		{
			"recovered object does not implement stackTracer",
			"/",
			func(w http.ResponseWriter, r *http.Request) {
				var i interface{}
				_ = i.(int) // cause an interface conversion panic
			},
			// LogUnstructuredStack generates the stack trace, so its name
			// ought to be present within it.
			"LogUnstructuredStack",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			// Configure a test router with recovery middleware.
			router := mux.NewRouter()
			logger := new(mockStackLogger)
			router.Use(Recovery(logger))
			router.HandleFunc(tc.path, tc.hf)

			// Make the request to the test path.
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tc.path, nil)
			router.ServeHTTP(rec, req)

			// Evaluate output
			if gotCode := rec.Result().StatusCode; gotCode != http.StatusInternalServerError {
				t.Errorf("want status code %d, got %d", http.StatusInternalServerError, gotCode)
			}

			output := logger.String()
			if !strings.Contains(output, tc.wantOutputToContain) {
				t.Errorf("want stack trace to contain %q, got trace %q",
					tc.wantOutputToContain, output)
			}
		})
	}
}
