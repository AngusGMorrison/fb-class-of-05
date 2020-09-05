package routing

import (
	"angusgmorrison/fb05/internal/app/templates"
	"bytes"
	"html/template"

	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type currentUser struct {
	LoggedIn bool
}

func getCurrentUser() *currentUser {
	return nil
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data := struct{ User *currentUser }{}
		tmpl, err := templates.Get("homepage")
		if err != nil {
			// Named templates must be initialized and accessible; this
			// should not happen.
			panic(err)
		}

		if err = writeTemplate(w, tmpl, data); err != nil {
			log.Error().Err(err)
		}
	default:
		http.NotFound(w, r)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tmpl, err := templates.Get("login")
		if err != nil {
			panic(err)
		}

		err = writeTemplate(w, tmpl, nil)
		if err != nil {
			log.Error().Err(err)
		}
	case http.MethodPost:
	default:
		http.NotFound(w, r)
	}
}

func templateError(handlerName string, err error, w http.ResponseWriter) {
	log.Error().Err(errors.Wrap(err, handlerName))
	http.Error(w, http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)
}

func writeTemplate(w http.ResponseWriter, t *template.Template, data interface{}) (err error) {
	var b bytes.Buffer
	if err = t.Execute(&b, data); err != nil {
		panic(err)
	}

	_, err = b.WriteTo(w)
	if err != nil {
		log.Error().Err(err)
	}

	return
}
