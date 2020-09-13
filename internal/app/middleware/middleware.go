// Package middleware provides custom middleware implementations for
// FB05.
package middleware

import (
	"net/http"
)

type WrappedWriter struct {
	http.Flusher
	http.ResponseWriter
	status int
	size   int
}

func WrapResponseWriter(w http.ResponseWriter) *WrappedWriter {
	return &WrappedWriter{ResponseWriter: w}
}

// Status returns the status captured when the underlying
// http.ResponseWriter was written. Returns 0 until the status header
// is written.
func (ww *WrappedWriter) Status() int {
	return ww.status
}

func (ww *WrappedWriter) WriteHeader(code int) {
	ww.status = code
	ww.ResponseWriter.WriteHeader(code)
}

// Write writes bytes b to the underlying http.ResponseWriter,
// automatically calling WriteHeader with http.StatusOK if it hasn't
// already been called. Records the number of bytes written.
func (ww *WrappedWriter) Write(b []byte) (int, error) {
	if !ww.Written() {
		ww.WriteHeader(http.StatusOK)
	}

	size, err := ww.ResponseWriter.Write(b)
	ww.size += size
	return size, err
}

// Written reports whether the http.ResponseWriter has already been
// written to.
func (ww *WrappedWriter) Written() bool {
	return ww.status != 0
}

// Size returns the number of bytes written to the
// http.ResponseWriter.
func (ww *WrappedWriter) Size() int {
	return ww.size
}
