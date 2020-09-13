package routing

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
)

func writeTemplate(w io.Writer, t *template.Template, data interface{}) error {
	var b bytes.Buffer
	if err := t.Execute(&b, data); err != nil {
		return err
	}

	_, err := b.WriteTo(w)
	return err
}

func notAllowed(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
