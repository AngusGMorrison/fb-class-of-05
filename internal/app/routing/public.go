package routing

import (
	"html/template"
	"log"
	"net/http"
)

var tmpls map[string]*template.Template

func init() {
	tmpls = make(map[string]*template.Template)
	tmpls["homepage"] = template.Must(
		template.ParseFiles(
			sharedTmplDir+"/application.gohtml",
			publicTmplDir+"/homepage.gohtml",
		))
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	if err := tmpls["homepage"].Execute(w, nil); err != nil {
		log.Printf("homepageHandler: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}
