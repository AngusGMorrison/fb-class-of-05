package routing

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
)

func templatesDir() string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, "..", "..", "..", "web", "templates")
}

func testMethodNotAllowed(router *mux.Router, r *http.Request, t *testing.T) {
	t.Run(fmt.Sprintf("%s %s 404s", r.Method, r.URL.EscapedPath()), func(t *testing.T) {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, r)
		if gotCode := rec.Code; gotCode != http.StatusMethodNotAllowed {
			t.Errorf("want status code %d, got %d", http.StatusMethodNotAllowed, gotCode)

		}
	})
}
