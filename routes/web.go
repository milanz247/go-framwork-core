package routes

import (
	"net/http"

	"ticket-sync-mvc/app/handlers"
	"ticket-sync-mvc/framework/dev"
	"ticket-sync-mvc/framework/env"
	"ticket-sync-mvc/framework/http/middleware"
	"ticket-sync-mvc/public"
)

func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/", handlers.API)
	mux.Handle("/static/", http.StripPrefix("/static/", public.Handler()))
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) })
	if env.Get("APP_ENV", "production") == "development" {
		reload := dev.NewReload("resources/views", "public/assets")
		mux.HandleFunc("GET /dev/version", reload.Version)
		mux.HandleFunc("GET /dev/refresh.js", reload.Script)
	}
	mux.HandleFunc("/", handlers.HTML)
	return middleware.Protect(mux)
}
