package middleware

import (
	"net/http"

	"github.com/pkg/errors"
)

// stackTracer implements a method to retrive the stack trace
// generated at the time the underlying error was created.
type stackTracer interface {
	StackTrace() errors.StackTrace
	Error() string
}

// Recovery logs stack traces from recovered errors. If the error
// implements stackTracer, the trace will be from the point of error
// creation. If not, the trace is generated within Recovery itself.
func Recovery(l logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hf := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if p := recover(); p != nil {
					w.WriteHeader(http.StatusInternalServerError)
					var err error
					var ok bool
					if err, ok = p.(stackTracer); !ok {
						err = errors.Errorf("%v", p)
					}
					l.Printf("%-8s %+v", info, err)
				}
			}()
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(hf)
	}
}
