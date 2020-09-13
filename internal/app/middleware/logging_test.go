package middleware

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"
)

func TestNewDebuggableLog(t *testing.T) {
	debug := []bool{true, false}
	for _, b := range debug {
		testLogger := log.New(nil, "", 0)
		wantLogger := &DebuggableLog{testLogger, b}
		gotLogger := NewDebuggableLog(testLogger, b)
		if !reflect.DeepEqual(wantLogger, gotLogger) {
			t.Errorf("want *DebuggableLogger %+v, got %+v", wantLogger, gotLogger)
		}
	}
}

func TestPrintf(t *testing.T) {
	testCases := []struct {
		debug   bool
		format  string
		data    []interface{}
		wantLog string
	}{
		{true, "%s %s", []interface{}{"DEBUG", "test log"}, "DEBUG test log\n"},
		{false, "%s %s", []interface{}{"DEBUG", "test log"}, ""},
		{true, "%s %s", []interface{}{"INFO", "test log"}, "INFO test log\n"},
		{false, "%s %s", []interface{}{"INFO", "test log"}, "INFO test log\n"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("debug=%t level=%s", tc.debug, tc.data[0]), func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewDebuggableLog(log.New(&buf, "", 0), tc.debug)
			logger.Printf(tc.format, tc.data...)
			gotLog := buf.String()
			if gotLog != tc.wantLog {
				t.Errorf("want log entry %q, got %q", tc.wantLog, gotLog)
			}
		})
	}
}

func TestLogging(t *testing.T) {
	okHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	errorHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}

	testCases := []struct {
		debug          bool
		path           string
		wantLogToMatch string
	}{
		{true, "/", `^DEBUG\s+GET /, status 200, \d+ µs\n$`},
		{false, "/", ""},
		{true, "/error", `^INFO\s+GET /error, status 500, \d+ µs\n$`},
		{false, "/error", `^INFO\s+GET /error, status 500, \d+ µs\n$`},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("debug=%t, path=%s", tc.debug, tc.path), func(t *testing.T) {
			var buf bytes.Buffer
			dl := NewDebuggableLog(log.New(&buf, "", 0), tc.debug)
			withLogging := Logging(dl)
			mux := http.NewServeMux()
			if tc.path == "/" {
				mux.Handle("/", withLogging(http.HandlerFunc(okHandler)))
			} else {
				mux.Handle("/error", withLogging(http.HandlerFunc(errorHandler)))
			}

			rec := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tc.path, nil)
			mux.ServeHTTP(rec, req)

			rgx := regexp.MustCompile(tc.wantLogToMatch)
			gotLog := buf.String()
			if !rgx.MatchString(gotLog) {
				t.Errorf("wanted log to match %q, got %q", tc.wantLogToMatch, gotLog)
			}
		})
	}
}
