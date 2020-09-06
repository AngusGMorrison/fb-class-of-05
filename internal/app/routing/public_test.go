package routing

import (
	"angusgmorrison/fb05/internal/app/templates"
	"bytes"
	"fmt"
	"html/template"
	"io"
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
			t.Errorf("got status code %d, want %d", status, http.StatusOK)
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
				"unexpected body\n\tgot: %s\n\twant: %s",
				gotBody, wantBody)
		}
	})

	notAllowed := []string{http.MethodConnect, http.MethodDelete, http.MethodHead,
		http.MethodOptions, http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodTrace}
	for _, method := range notAllowed {
		req, _ := http.NewRequest(method, "/", nil)
		testMethodNotAllowed(router, req, t)
	}
}

func TestLoginHandler(t *testing.T) {
	if err := templates.Initialize(templatesDir()); err != nil {
		t.Fatalf("templates failed to initialize: %v", err)
	}
	router := Router(nil)

	tmpl, err := templates.Get("login")
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		method   string
		body     io.Reader
		wantCode int
		tmpl     *template.Template
		data     interface{}
	}{
		{http.MethodGet, nil, http.StatusOK, tmpl, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, "/login", tc.body)
			router.ServeHTTP(rec, req)

			if gotCode := rec.Code; gotCode != tc.wantCode {
				t.Errorf("got status code %d, want %d", gotCode, tc.wantCode)
			}

			if tc.tmpl != nil {
				var b bytes.Buffer
				err := tc.tmpl.Execute(&b, tc.data)
				if err != nil {
					t.Fatal(err)
				}

				if gotBody, wantBody := rec.Body.String(), b.String(); gotBody != wantBody {
					t.Errorf(
						"unexpected body\n\tgot: %s\n\twant: %s",
						gotBody, wantBody)
				}
			}
		})
	}

	notAllowed := []string{http.MethodConnect, http.MethodDelete, http.MethodHead,
		http.MethodOptions, http.MethodPatch, http.MethodPut, http.MethodTrace}

	for _, method := range notAllowed {
		req, _ := http.NewRequest(method, "/login", nil)
		testMethodNotAllowed(router, req, t)
	}
}

func templatesDir() string {
	cwd, _ := os.Getwd()
	return filepath.Join(cwd, "..", "templates")
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
