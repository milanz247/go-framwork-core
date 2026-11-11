package handlers

import (
	"crypto/hmac"
	"log"
	"net/http"
	"strings"

	"ticket-sync-mvc/app/auth"
	"ticket-sync-mvc/framework/env"
	"ticket-sync-mvc/framework/http/response"
	"ticket-sync-mvc/framework/support"
)

func API(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("request failed: %v", e)
			response.Fail(w, 500, "Request failed. Check server logs.")
		}
	}()
	path := strings.TrimPrefix(r.URL.Path, "/api/")
	if r.Method != "GET" {
		origin := r.Header.Get("Origin")
		if origin != "" {
			scheme := "http://"
			if r.TLS != nil {
				scheme = "https://"
			}
			if origin != scheme+r.Host && origin != env.Get("PUBLIC_ORIGIN", "") {
				response.Fail(w, 403, "Invalid request origin")
				return
			}
		}
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			response.Fail(w, 415, "JSON required")
			return
		}
	}
	if path == "login" || path == "register" {
		handleAuth(w, r, path)
		return
	}

	u := auth.SessionUser(r)
	if u == nil {
		response.Fail(w, 401, "Sign in to continue")
		return
	}
	if r.Method != "GET" && !hmac.Equal([]byte(r.Header.Get("X-CSRF-Token")), []byte(support.Str(u["csrf"]))) {
		response.Fail(w, 403, "Refresh and try again")
		return
	}
	switch {
	case path == "state" && r.Method == "GET":
		response.JSON(w, support.M{"user": u})
	case path == "logout" && r.Method == "POST":
		handleLogout(w, r)
	default:
		response.Fail(w, 404, "Not found")
	}
}
