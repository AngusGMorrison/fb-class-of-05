// Package routing describes the router and handlers for FB05.
package routing

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	templateDir   = "internal/app/templates"
	sharedTmplDir = templateDir + "/shared"
	publicTmplDir = templateDir + "/public"
)

// Router returns a mux that handles routing for the entire
// application.
func Router(mw []mux.MiddlewareFunc) *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/", homepageHandler)

	for _, m := range mw {
		router.Use(m)
	}

	return router
}
