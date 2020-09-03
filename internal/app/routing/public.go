package routing

import (
	"html/template"
	"log"
	"net/http"
)

type currentUser struct {
	LoggedIn bool
}

func getCurrentUser() *currentUser {
	return nil
}

var tmpls map[string]*template.Template

func init() {
	tmpls = make(map[string]*template.Template)
	tmpls["homepage"] = template.Must(
		template.ParseFiles(
			sharedTmplDir+"/application.gohtml",
			sharedTmplDir+"/banner_nav.gohtml",
			sharedTmplDir+"/sidebar_nav.gohtml",
			publicTmplDir+"/homepage.gohtml",
		))
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	data := struct{ User *currentUser }{}
	if err := tmpls["homepage"].Execute(w, data); err != nil {
		log.Printf("homepageHandler: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
}
