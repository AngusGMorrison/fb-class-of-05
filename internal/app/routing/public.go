package routing

import (
	"angusgmorrison/fb05/internal/app/templates"
	"net/http"
)

type currentUser struct {
	LoggedIn bool
}

func getCurrentUser() *currentUser {
	return nil
}

func homepageHandler(w http.ResponseWriter, r *http.Request, log logger) {
	switch r.Method {
	case http.MethodGet:
		data := struct{ User *currentUser }{}
		tmpl, err := templates.Get("homepage")
		if err != nil {
			panic(err) // should not happen
		}

		if err = writeTemplate(w, tmpl, data); err != nil {
			log.Printf("%-8s %s: %v", "INFO", "homepageHandler", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
	default:
		notAllowed(w, r)
		return
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request, log logger) {
	switch r.Method {
	case http.MethodGet:
		tmpl, err := templates.Get("login")
		if err != nil {
			panic(err) // should not happen
		}

		err = writeTemplate(w, tmpl, nil)
		if err != nil {
			log.Printf("%-8s %s: %v", "INFO", "loginHandler", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
	default:
		notAllowed(w, r)
		return
	}
}
