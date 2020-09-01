package middleware

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// A StackLogger implements methods to log a StackTracer's stack trace
// as structured JSON or to generate new stack trace and log it as
// an unstructured string.
type StackLogger interface {
	LogStructuredStack(err StackTracer)
	LogUnstructuredStack(err error)
}

// A StackTracer implements a method to retrive the stack trace
// generated at the time the underlying error was created.
type StackTracer interface {
	StackTrace() errors.StackTrace
}

// Recovery accepts a StackLogger and returns a function that wraps
// an http.Handler in panic-recovery middleware which sets
// StatusInternalServerError and logs either a stuctured or
// unstructured stack trace depending on the concrete type recovered
// from the panic.
func Recovery(logger StackLogger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		hf := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if p := recover(); p != nil {
					w.WriteHeader(http.StatusInternalServerError)
					if tracer, ok := p.(StackTracer); ok {
						// We have access to stack frames generated at
						// the point of error creation, so log them as
						// JSON.
						logger.LogStructuredStack(tracer)
					} else {
						// At best we can generate a string of the
						// stack trace at the current execution point.
						err := fmt.Errorf("%v", p)
						logger.LogUnstructuredStack(err)
					}
				}
			}()
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(hf)
	}
}
