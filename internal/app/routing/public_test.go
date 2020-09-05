package routing

import (
	"angusgmorrison/fb05/internal/app/templates"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gorilla/mux"
)

func TestHomepageHandler(t *testing.T) {
	if err := templates.Initialize(templatesDir()); err != nil {
		t.Fatalf("templates failed to initialize: %v", err)
	}
	router := Router(nil)

	t.Run("method GET", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		router.ServeHTTP(rec, req)

		if status := rec.Code; status != http.StatusOK {
			t.Errorf("homepageHandler returned status code %d, want %d", status, http.StatusOK)
		}

		tmpl, err := templates.Get("homepage")
		if err != nil {
			t.Fatal(err)
		}
		var b bytes.Buffer
		err = tmpl.Execute(&b, nil)
		if err != nil {
			t.Fatal(err)
		}
		if gotBody, wantBody := rec.Body.String(), b.String(); gotBody != wantBody {
			t.Errorf(
				"homepageHandler returned unexpected body\n\tgot: %s\n\twant: %s",
				gotBody, wantBody)
		}
	})

	should404 := []string{http.MethodConnect, http.MethodDelete, http.MethodHead,
		http.MethodOptions, http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodTrace}
	for _, method := range should404 {
		req, _ := http.NewRequest(method, "/", nil)
		test404(router, req, t)
	}
}

func templatesDir() string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, "..", "templates")
}

func test404(router *mux.Router, r *http.Request, t *testing.T) {
	t.Run(fmt.Sprintf("%s %s 404s", r.Method, r.URL.EscapedPath()), func(t *testing.T) {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, r)
		if gotCode := rec.Code; gotCode != http.StatusNotFound {
			t.Errorf("want status code %d, got %d", http.StatusNotFound, gotCode)

		}
	})
}
