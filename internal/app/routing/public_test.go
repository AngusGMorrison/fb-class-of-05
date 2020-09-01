package routing

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHomepageHandler(t *testing.T) {
	router := Router(nil)
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("homepageHandler returned status code %d, want %d", status, http.StatusOK)
	}

	wantBody := "Hello, World!"
	if gotBody := rec.Body.String(); gotBody != wantBody {
		t.Errorf(
			"homepageHandler returned unexpected body\n\tgot: %s\n\twant: %s",
			gotBody,
			wantBody,
		)
	}
}
