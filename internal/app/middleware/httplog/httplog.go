// Package httplog provides a mechanism for logging HTTP events by
// wrapping a *zerolog.Logger while satisfying middleware.HTTPLogger.
package httplog

import (
	"angusgmorrison/fb05/internal/app/middleware"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	*zerolog.Logger
}

func NewLogger(logger *zerolog.Logger) *Logger {
	return &Logger{logger}
}

// Log logs details of the HTTP request and response status using the
// Logger's underlying *zerolog.Logger.
func (l *Logger) Log(ww *middleware.WrappedWriter, r *http.Request, d time.Duration) {
	var event *zerolog.Event
	status := ww.Status()
	if status == http.StatusInternalServerError {
		event = l.Error()
	} else {
		event = l.Info()
	}

	event.Str("method", r.Method).
		Str("path", r.URL.EscapedPath()).
		Int("status", status).
		Dur("duration", d).
		Send()
}
