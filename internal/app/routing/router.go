// Package routing describes the router and handlers for FB05.
package routing

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Router returns a mux that handles routing for the entire
// application.
func Router(mw []mux.MiddlewareFunc) *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/login", loginHandler).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/", homepageHandler).Methods(http.MethodGet)

	for _, m := range mw {
		router.Use(m)
	}

	return router
}
