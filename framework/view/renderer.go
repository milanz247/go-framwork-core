// Package view loads nested templates and renders buffered HTML responses.
package view

import (
	"bytes"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
)

type Renderer struct{ templates *template.Template }

func New(source fs.FS, root string, functions template.FuncMap) (*Renderer, error) {
	paths := []string{}
	err := fs.WalkDir(source, root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && strings.HasSuffix(path, ".html") {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	templates, err := template.New("views").Funcs(helpers()).Funcs(functions).ParseFS(source, paths...)
	if err != nil {
		return nil, err
	}
	return &Renderer{templates: templates}, nil
}

func (renderer *Renderer) Render(w http.ResponseWriter, status int, name string, data any) {
	var body bytes.Buffer
	if err := renderer.templates.ExecuteTemplate(&body, name, data); err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	w.Write(body.Bytes())
}
