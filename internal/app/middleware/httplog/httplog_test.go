package httplog

import (
	"angusgmorrison/fb05/internal/app/middleware"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
)

const env = "test"

func TestNewLogger(t *testing.T) {
	zerologger := &log.Logger
	wantLogger := &Logger{zerologger, env}
	gotLogger := NewLogger(zerologger, env)
	if !reflect.DeepEqual(wantLogger, gotLogger) {
		t.Errorf("Expected returned logger to match %+v, got %+v", wantLogger, gotLogger)
	}
}

func TestLoggerImplementsHTTPLogger(t *testing.T) {
	var _ middleware.HTTPLogger = NewLogger(&log.Logger, env)
}

func newTestLogger(env string) (*Logger, *bytes.Buffer) {
	var output bytes.Buffer
	zeroLogger := log.Output(&output)
	return NewLogger(&zeroLogger, env), &output
}

func TestLog(t *testing.T) {
	testCases := []struct {
		level    string
		method   string
		path     string
		status   int
		duration time.Duration
	}{
		{"info", "GET", "/", http.StatusOK, 1},
		{"error", "GET", "/", http.StatusInternalServerError, 1},
	}

	for _, tc := range testCases {
		logger, output := newTestLogger(env)
		ww := middleware.WrapResponseWriter(httptest.NewRecorder())
		ww.WriteHeader(tc.status)
		req, _ := http.NewRequest("GET", "/", nil)
		logger.Log(ww, req, tc.duration*time.Millisecond)

		gotLog := output.String()
		wantLog := fmt.Sprintf(
			`{"level":%q,"method":%q,"path":%q,"status":%d,"duration":%d,"time":"[\d\-T:\+]{25}"}`,
			tc.level, tc.method, tc.path, tc.status, tc.duration)
		rx := regexp.MustCompile(wantLog)
		if !rx.MatchString(gotLog) {
			t.Errorf("want log entry %q, got %q", wantLog, gotLog)
		}
	}
}
