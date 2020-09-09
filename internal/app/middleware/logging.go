package middleware

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	info  = "INFO"
	debug = "DEBUG"
)

type logger interface {
	Printf(format string, v ...interface{})
}

// DebuggableLog signals whether to log in debug mode to its embedded
// logger.
type DebuggableLog struct {
	logger
	debug bool
}

func NewDebuggableLog(l logger, debug bool) *DebuggableLog {
	return &DebuggableLog{l, debug}
}

// Printf checks if debug mode is enabled before calling
// logger.Printf. If not, it attempts to read the current logging
// level from v[0] and ignores submissions of level debug.
func (dl *DebuggableLog) Printf(format string, v ...interface{}) {
	if dl.debug {
		dl.logger.Printf(format, v...)
		return
	}

	level, ok := v[0].(string)
	if !ok || level != debug {
		dl.logger.Printf(format, v...)
	}
}

// Logging wraps an HTTP handler to log request details.
func Logging(l logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		hf := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := WrapResponseWriter(w)
			next.ServeHTTP(ww, r)

			var level string
			if ww.Status() == http.StatusInternalServerError {
				level = info
			} else {
				level = debug
			}
			l.Printf("%-8s %s %s, status %d, %d µs\n",
				level, r.Method, r.URL.EscapedPath(), ww.Status(), time.Since(start).Microseconds())
		}

		return http.HandlerFunc(hf)
	}
}
