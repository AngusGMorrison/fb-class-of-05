package routing

import (
	"angusgmorrison/fb05/internal/app/templates"
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kylelemons/godebug/diff"
)

func TestHomepageHandler(t *testing.T) {
	if err := templates.Initialize(templatesDir()); err != nil {
		t.Fatalf("templates failed to initialize: %v", err)
	}

	var buf bytes.Buffer
	handler := &loggingHandler{log.New(&buf, "", 0), homepageHandler}
	mux := http.NewServeMux()
	mux.Handle("/", handler)

	t.Run(http.MethodGet, func(t *testing.T) {
		getTmpl, err := templates.Get("homepage")
		if err != nil {
			t.Fatal(err)
		}

		var buf bytes.Buffer
		err = writeTemplate(&buf, getTmpl, nil)
		if err != nil {
			t.Fatal(err)
		}
		wantBody := buf.String()

		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("want status code %d, got %d", http.StatusOK, rec.Code)
		}

		if gotBody := rec.Body.String(); gotBody != wantBody {
			d := diff.Diff(wantBody, gotBody)
			t.Errorf("received body doesn't match expected body:\n%s", d)
		}
	})

	notAllowed := []string{http.MethodConnect, http.MethodDelete, http.MethodHead,
		http.MethodOptions, http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodTrace}

	for _, method := range notAllowed {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(method, "/", nil)
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("want status code %d, got %d", http.StatusMethodNotAllowed, rec.Code)
		}
	}
}

func TestLoginHandler(t *testing.T) {
	if err := templates.Initialize(templatesDir()); err != nil {
		t.Fatalf("templates failed to initialize: %v", err)
	}

	var buf bytes.Buffer
	handler := &loggingHandler{log.New(&buf, "", 0), loginHandler}
	mux := http.NewServeMux()
	mux.Handle("/login", handler)

	getTmpl, err := templates.Get("login")
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
		{http.MethodGet, nil, http.StatusOK, getTmpl, nil},
		// TODO: POST method test
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, "/login", tc.body)
			mux.ServeHTTP(rec, req)

			if gotCode := rec.Code; gotCode != tc.wantCode {
				t.Errorf("got status code %d, want %d", gotCode, tc.wantCode)
			}

			var buf bytes.Buffer
			err := writeTemplate(&buf, tc.tmpl, tc.data)
			if err != nil {
				t.Fatal(err)
			}

			if gotBody, wantBody := rec.Body.String(), buf.String(); gotBody != wantBody {
				d := diff.Diff(wantBody, gotBody)
				t.Errorf("received body doesn't match expected body:\n%s", d)
			}
		})
	}

	notAllowed := []string{http.MethodConnect, http.MethodDelete, http.MethodHead,
		http.MethodOptions, http.MethodPatch, http.MethodPut, http.MethodTrace}

	for _, method := range notAllowed {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(method, "/login", nil)
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("want status code %d, got %d", http.StatusMethodNotAllowed, rec.Code)
		}
	}
}
