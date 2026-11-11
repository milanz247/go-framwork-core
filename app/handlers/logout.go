package handlers

import (
	"net/http"

	"ticket-sync-mvc/app/auth"
	"ticket-sync-mvc/framework/env"
	"ticket-sync-mvc/framework/http/response"
	"ticket-sync-mvc/framework/support"
)

func handleLogout(w http.ResponseWriter, r *http.Request) {
	auth.EndSession(r)
	http.SetCookie(w, &http.Cookie{Name: "ticket_session", Path: "/", MaxAge: -1, HttpOnly: true, Secure: env.Get("COOKIE_SECURE", "0") == "1", SameSite: http.SameSiteLaxMode})
	if r.URL.Path == "/logout" {
		redirect(w, r, "/login")
		return
	}
	response.JSON(w, support.M{"ok": true})
}
