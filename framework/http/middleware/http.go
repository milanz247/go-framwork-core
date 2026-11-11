package middleware

import (
	"log"
	"net/http"
	"time"

	"ticket-sync-mvc/framework/env"
)

func Protect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
			if e := recover(); e != nil {
				log.Printf("request panic: %v", e)
				http.Error(w, "Request failed. Check server logs.", 500)
			}
		}()
		w.Header().Set("Cache-Control", "private, no-store")
		w.Header().Add("Vary", "HX-Request")
		w.Header().Add("Vary", "HX-History-Restore-Request")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "same-origin")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; connect-src 'self'; object-src 'none'; base-uri 'self'; frame-ancestors 'none'; form-action 'self'")
		if r.Method != "GET" && r.Method != "HEAD" {
			origin := r.Header.Get("Origin")
			scheme := "http://"
			if r.TLS != nil {
				scheme = "https://"
			}
			allowed := env.Get("PUBLIC_ORIGIN", scheme+r.Host)
			if (origin != "" && origin != allowed) || r.Header.Get("Sec-Fetch-Site") == "cross-site" {
				http.Error(w, "Invalid request origin", 403)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
