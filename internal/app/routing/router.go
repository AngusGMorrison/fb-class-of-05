// Package routing describes the router and handlers for FB05.
package routing

import (
	"net/http"

	"github.com/gorilla/mux"
)

const staticDir = "web/static"

// Router returns a mux that handles routing for the entire
// application.
func Router(log logger, mw []func(http.Handler) http.Handler) *mux.Router {
	router := mux.NewRouter()
	withLogger := loggingHandlerFactory(log)

	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))
	router.Handle("/login", withLogger(loginHandler)).Methods(http.MethodGet, http.MethodPost)
	router.Handle("/", withLogger(homepageHandler)).Methods(http.MethodGet)

	for _, m := range mw {
		router.Use(m)
	}

	return router
}
