// Package templates handles template parsing and access for FB05.
package templates

import (
	"html/template"
	"path/filepath"

	"github.com/pkg/errors"
)

type templateInstruction struct {
	name         string
	baseTemplate string
	files        []string
}

var tmpls map[string]*template.Template

// Initialize loads all required templates into memory. It is NOT
// thread-safe.
// TODO: Convert to concurrent, non-blocking cache.
func Initialize(templateDir string) error {
	var (
		sharedDir = filepath.Join(templateDir, "shared")
		publicDir = filepath.Join(templateDir, "public")
	)
	templateInstructions := []templateInstruction{
		{
			"root",
			"",
			[]string{filepath.Join(sharedDir, "application.gohtml"),
				filepath.Join(sharedDir, "banner_nav.gohtml")},
		},
		{"primary_sidebar_base", "root", []string{filepath.Join(sharedDir, "sidebar_nav.gohtml")}},
		{"homepage", "primary_sidebar_base", []string{filepath.Join(publicDir, "homepage.gohtml")}},
		{"login", "primary_sidebar_base", []string{filepath.Join(publicDir, "login.gohtml")}},
	}

	var err error
	tmpls = make(map[string]*template.Template)
	for _, ti := range templateInstructions {
		if ti.baseTemplate == "" {
			tmpls[ti.name], err = template.ParseFiles(ti.files...)
			if err != nil {
				return errors.Errorf("parsing template %q: %v", ti.name, err)
			}
			continue
		}

		base, ok := tmpls[ti.baseTemplate]
		if !ok {
			return errors.Errorf("base template %q does not exist", ti.name)
		}
		t, err := base.Clone()
		if err != nil {
			return errors.WithStack(err)
		}
		tmpls[ti.name], err = t.ParseFiles(ti.files...)
		if err != nil {
			return errors.Errorf("parsing template %q: %v", ti.name, err)
		}
	}

	return nil
}

// Get returns a named template from the cache of parsed templates.
func Get(name string) (*template.Template, error) {
	tmpl, ok := tmpls[name]
	if !ok {
		return nil, errors.Errorf("template %q not found", name)
	}
	clone, err := tmpl.Clone()
	if err != nil {
		return nil, errors.Wrapf(err, "cloning template %q", name)
	}
	return clone, nil
}
