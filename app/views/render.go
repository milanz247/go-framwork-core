// Package views selects the application layout for full pages and HTMX requests.
package views

import (
	"net/http"

	"ticket-sync-mvc/framework/support"
	"ticket-sync-mvc/framework/view"
	"ticket-sync-mvc/resources"
)

var renderer = newRenderer()

func newRenderer() *view.Renderer {
	renderer, err := view.New(resources.Views, "views", nil)
	if err != nil {
		panic(err)
	}
	return renderer
}

func Render(w http.ResponseWriter, r *http.Request, status int, data support.M) {
	layout := "layout/base"
	if r.Header.Get("HX-Request") == "true" && r.Header.Get("HX-History-Restore-Request") != "true" {
		layout = "layout/app"
		if data["Auth"] == true {
			layout = "layout/auth"
		}
	}
	renderer.Render(w, status, layout, data)
}
