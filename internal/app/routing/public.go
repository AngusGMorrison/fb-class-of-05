package routing

import (
	"html/template"
	"log"
	"net/http"
)

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("internal/app/templates/shared/application.gohtml",
		"internal/app/templates/public/homepage.gohtml")
	if err != nil {
		log.Printf("homepageHandler: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	if err = t.Execute(w, nil); err != nil {
		log.Printf("homepageHandler: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}
