package middleware

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecovery(t *testing.T) {
	tracingPanicHandler := func(w http.ResponseWriter, r *http.Request) {
		panic(errors.New("panic")) // implements stackTracer
	}
	nontracingPanicHandler := func(w http.ResponseWriter, r *http.Request) {
		var i interface{}
		_ = i.(int) // panic generated doesn't implement stackTracer
	}

	testCases := []struct {
		path               string
		handlerFn          http.HandlerFunc
		wantCode           int
		wantTraceToContain string
	}{
		{"/stackTracer", tracingPanicHandler, http.StatusInternalServerError, "TestRecovery"},
		{"/notStackTracer", nontracingPanicHandler, http.StatusInternalServerError, "TestRecovery"},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			var buf bytes.Buffer
			testLogger := log.New(&buf, "", 0)
			withRecovery := Recovery(testLogger)
			mux := http.NewServeMux()
			mux.Handle(tc.path, withRecovery(http.HandlerFunc(tc.handlerFn)))

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tc.path, nil)
			mux.ServeHTTP(rec, req)

			if rec.Code != tc.wantCode {
				t.Errorf("want status code %d, got %d", tc.wantCode, rec.Code)
			}
			if gotLog := buf.String(); !strings.Contains(gotLog, tc.wantTraceToContain) {
				t.Errorf("want logged stack trace to contain %q, got trace %q",
					tc.wantTraceToContain, gotLog)
			}
		})
	}
}
