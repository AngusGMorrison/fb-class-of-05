// Package stacklog facilitates logging stack traces as structured
// JSON and unstructured strings by embedding a *zerolog.Logger in a
// Logger struct that implements internal/app/middleware.Stacklogger.
package stacklog

import (
	"angusgmorrison/fb05/internal/app/middleware"
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// prettyStackWriter is the location stack traces are printed to in
// a development environment.
var prettyStackWriter io.Writer = os.Stderr

type Logger struct {
	*zerolog.Logger
	env string
}

func NewLogger(logger *zerolog.Logger, env string) *Logger {
	return &Logger{logger, env}
}

// LogStructuredStack logs a JSON stack generated at point of error
// creation. To do so, the error must implement
// middleware.StackTracer. This will handle logging for panics written
// by the authors.
func (l *Logger) LogStructuredStack(tracer middleware.StackTracer) {
	err, ok := tracer.(error)
	if !ok {
		panic(errors.New(
			fmt.Sprintf("StackTracer should have concrete type 'error', got %T", tracer)))
	}

	l.Error().
		Stack().
		Err(err).
		Send()

	if l.env == "development" {
		// Print formatted stack for readability.
		fmt.Fprintln(prettyStackWriter, "\nDEVELOPMENT ONLY:")
		for _, f := range tracer.StackTrace() {
			fmt.Fprintf(prettyStackWriter, "%+s:%d\n", f, f)
		}
	}
}

// LogUnstructuredStack logs a string stack generated at the current
// point of execution (i.e., within logUnstructuredStack). This
// handles logging for panics that aren't known to be possible ahead
// of time.
func (l *Logger) LogUnstructuredStack(err error) {
	stack := string(debug.Stack())
	l.Error().
		Err(err).
		Str("stack", stack).
		Send()

	if l.env == "development" {
		// Print formatted stack for readability.
		fmt.Fprintln(prettyStackWriter, "\nDEVELOPMENT ONLY:")
		fmt.Fprintln(prettyStackWriter, stack)
	}
}
