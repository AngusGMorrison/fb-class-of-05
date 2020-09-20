// Package templates handles template parsing and access for FB05
// via a concurrent, non-blocking cache.
package templates

import (
	"angusgmorrison/fb05/pkg/env"
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type result struct {
	tmpl *template.Template
	err  error
}

func Get(name string) (*template.Template, error) {
	if requestStream == nil {
		return nil, errors.New("template cache has not been initialized")
	}
	resultStream := make(chan *result)

	timeout := time.Duration(env.Get("GET_TIMEOUT").(int)) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Add timeout/cancellation
	req := &request{ctx, name, resultStream}
	requestStream <- req
	res := <-resultStream
	if res.err != nil {
		return nil, res.err
	}
	return res.tmpl, nil
}

type request struct {
	ctx          context.Context
	name         string
	resultStream chan<- *result
}

var requestStream chan<- *request

func InitCache(done <-chan struct{}, templateDir string) {
	requestStream = cacheTemplates(done, templateDir)
}

// templateSchematic holds instructions for building a named template
// from a base template and constituent files.
type templateSchematic struct {
	baseTemplate string
	files        []string
}

type cacheEntry struct {
	ready chan struct{}
	tmpl  *template.Template
	err   error
}

func cacheTemplates(done <-chan struct{}, templateDir string) chan<- *request {
	var (
		sharedDir = filepath.Join(templateDir, "shared")
		publicDir = filepath.Join(templateDir, "public")
	)

	schematics := map[string]*templateSchematic{
		"root": {
			"",
			[]string{filepath.Join(sharedDir, "application.gohtml"),
				filepath.Join(sharedDir, "banner_nav.gohtml")},
		},
		"primary_sidebar_base": {"root", []string{filepath.Join(sharedDir, "sidebar_nav.gohtml")}},
		"homepage":             {"primary_sidebar_base", []string{filepath.Join(publicDir, "homepage.gohtml")}},
		"login":                {"primary_sidebar_base", []string{filepath.Join(publicDir, "login.gohtml")}},
	}

	templates := make(map[string]*cacheEntry)
	requestStream := make(chan *request)
	go func() {
		defer close(requestStream)

		for {
			select {
			case <-done:
				return
			case req := <-requestStream:
				select {
				case <-req.ctx.Done():
					req.resultStream <- &result{err: req.ctx.Err()}
				default:
				}

				entry := templates[req.name]
				if entry == nil {
					schematic := schematics[req.name]
					if schematic == nil {
						req.resultStream <- &result{
							err: errors.New(
								fmt.Sprintf("requested schematic %q does not exist", req.name)),
						}
						continue
					}

					entry = &cacheEntry{ready: make(chan struct{})}
					templates[req.name] = entry
					go entry.parse(schematic)
				}
				go entry.deliver(req.resultStream)
			}
		}
	}()

	return requestStream
}

func (ce *cacheEntry) parse(s *templateSchematic) {
	defer close(ce.ready)

	var tmpl *template.Template
	var err error
	if s.baseTemplate == "" {
		tmpl, err = template.ParseFiles(s.files...)
	} else {
		base, err := Get(s.baseTemplate)
		if err != nil {
			ce.err = err
			return
		}
		tmpl, err = base.ParseFiles(s.files...)
	}
	if err != nil {
		ce.err = err
		return
	}

	ce.tmpl = tmpl
	return
}

func (ce *cacheEntry) deliver(resultStream chan<- *result) {
	<-ce.ready // Make this preemptible

	if ce.err != nil {
		resultStream <- &result{err: ce.err}
		return
	}

	clone, err := ce.tmpl.Clone()
	if err != nil {
		resultStream <- &result{err: ce.err}
		return
	}
	resultStream <- &result{tmpl: clone}
	return
}
