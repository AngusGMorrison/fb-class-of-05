package middleware

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// HTTPLogger provides a simple interface through which HTTP events
// can be recorded.
type HTTPLogger interface {
	Log(ww *WrappedWriter, r *http.Request, duration time.Duration)
}

// Logging accepts a preconfigured *zerolog.Logger and
// returns a function that returns the middleware-wrapped http.Handler
// when called.
func Logging(logger HTTPLogger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		hf := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := WrapResponseWriter(w)
			next.ServeHTTP(ww, r)
			logger.Log(ww, r, time.Since(start))
		}

		return http.HandlerFunc(hf)
	}
}
