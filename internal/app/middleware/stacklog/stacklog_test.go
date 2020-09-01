package stacklog

import (
	"angusgmorrison/fb05/internal/app/middleware"
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

const env = "test"

func TestNewLogger(t *testing.T) {
	zeroLogger := log.Logger
	wantLogger := &Logger{&zeroLogger, env}
	gotLogger := NewLogger(&zeroLogger, env)
	if !reflect.DeepEqual(wantLogger, gotLogger) {
		t.Errorf("want *Logger %+v, got *Logger %+v", wantLogger, gotLogger)
	}
}

func TestLoggerImplementsStackLogger(t *testing.T) {
	var _ middleware.StackLogger = NewLogger(&log.Logger, env)
}

// errorlessStackTracer is used to test LogStructuredStack with a
// middleware.StackTracer that doesn't have concrete type 'error'.
type errorlessStackTracer struct {
	msg string
}

func (e errorlessStackTracer) StackTrace() errors.StackTrace {
	return errors.StackTrace{}
}

func newTestLogger(env string) (*Logger, *bytes.Buffer) {
	var output bytes.Buffer
	zeroLogger := log.Output(&output)
	return NewLogger(&zeroLogger, env), &output
}

func TestLogStructuredStack(t *testing.T) {
	// Enable stack trace logging for zerolog.
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	t.Run("logs stack traces in any environment", func(t *testing.T) {
		testLogger, output := newTestLogger(env)
		tracer := errors.New("test error").(middleware.StackTracer)
		testLogger.LogStructuredStack(tracer)

		logged := output.String()
		wantErrMsg := `"error":"test error"`
		if !strings.Contains(logged, wantErrMsg) {
			t.Errorf("want logged error to contain %q, got %q", wantErrMsg, logged)
		}
		wantStackPrefix := fmt.Sprintf("%q:[{%q:%q",
			"stack", "func", "TestLogStructuredStack.func1")
		if !strings.Contains(logged, wantStackPrefix) {
			t.Errorf("want logged stack trace to contain %q, got %q", wantStackPrefix, logged)
		}
	})

	t.Run("panics when StackTracer is not an error", func(t *testing.T) {
		defer func() {
			if p := recover(); p == nil {
				t.Errorf("Passing errorlessStackTracer didn't cause a panic.")
			}
		}()

		testLogger, _ := newTestLogger(env)
		testLogger.LogStructuredStack(errorlessStackTracer{"Should panic"})
	})

	t.Run("pretty prints stack traces in development", func(t *testing.T) {
		// Mock os.Stderr.
		oldPSW := prettyStackWriter
		prettyStackWriter = new(bytes.Buffer)
		defer func() {
			prettyStackWriter = oldPSW
		}()

		testLogger, _ := newTestLogger("development")
		tracer := errors.New("test error").(middleware.StackTracer)
		testLogger.LogStructuredStack(tracer)

		// Validate that output looks like a stack trace.
		output := prettyStackWriter.(*bytes.Buffer).String()
		if !strings.Contains(output, "stacklog_test.go") {
			t.Errorf("stack trace failed to print correctly, got:\n%q", output)
		}
	})
}

func TestLogUnstructuredStack(t *testing.T) {
	t.Run("logs stack traces in any environment", func(t *testing.T) {
		logger, output := newTestLogger(env)
		err := errors.New("test error")
		logger.LogUnstructuredStack(err)

		logged := output.String()
		wantErrMsg := `"error":"test error"`
		if !strings.Contains(logged, wantErrMsg) {
			t.Errorf("want logged error to contain %q, got %q", wantErrMsg, logged)
		}
		// LogUnstructuredStack generates the stack trace, so its name
		// ought to be present within it.
		wantStackToContain := "LogUnstructuredStack"
		if !strings.Contains(logged, wantStackToContain) {
			t.Errorf("want logged stack trace to contain %q, got %q", wantStackToContain, logged)
		}
	})

	t.Run("pretty prints stack traces in development", func(t *testing.T) {
		// Mock os.Stderr.
		oldPSW := prettyStackWriter
		prettyStackWriter = new(bytes.Buffer)
		defer func() {
			prettyStackWriter = oldPSW
		}()

		testLogger, _ := newTestLogger("development")
		testLogger.LogUnstructuredStack(errors.New("test error"))

		// Validate that output looks like a stack trace.
		output := prettyStackWriter.(*bytes.Buffer).String()
		if !strings.Contains(output, "stacklog_test.go") {
			t.Errorf("stack trace failed to print correctly, got:\n%q", output)
		}
	})
}
