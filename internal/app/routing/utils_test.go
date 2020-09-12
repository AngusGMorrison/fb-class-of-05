package routing

import (
	"bytes"
	"fmt"
	"html/template"
	"testing"
)

func TestWriteTemplate(t *testing.T) {
	rawTmpl := `
	<html>
		<body>
			<p>Test template</p>
				{{if .}}
				<p>Data!</p>
				{{end}}
		</body>
	</html>`

	tmpl := template.Must(template.New("test").Parse(rawTmpl))

	testData := []interface{}{true, nil}
	for _, data := range testData {
		t.Run(fmt.Sprintf("data=%v", data), func(t *testing.T) {
			var want bytes.Buffer
			err := tmpl.Execute(&want, data)
			if err != nil {
				t.Fatal(err)
			}

			var got bytes.Buffer
			err = writeTemplate(&got, tmpl, data)
			if err != nil {
				t.Fatal(err)
			}

			if gotStr, wantStr := got.String(), want.String(); gotStr != wantStr {
				t.Errorf("want %q, got %q", wantStr, gotStr)
			}

		})
	}
}
