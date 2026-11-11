package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"ticket-sync-mvc/app/repositories"
	"ticket-sync-mvc/framework/env"
	"ticket-sync-mvc/framework/support"
)

func SessionUser(r *http.Request) support.M {
	c, e := r.Cookie("ticket_session")
	if e != nil {
		return nil
	}
	key := sha256.Sum256([]byte(c.Value))
	return repositories.FindSession(hex.EncodeToString(key[:]))
}

func StartSession(w http.ResponseWriter, id int64) {
	token := support.Random()
	key := sha256.Sum256([]byte(token))
	repositories.SaveSession(hex.EncodeToString(key[:]), id, support.Random(), time.Now().Add(24*time.Hour).UTC().Format(time.RFC3339))
	http.SetCookie(w, &http.Cookie{Name: "ticket_session", Value: token, Path: "/", HttpOnly: true, Secure: env.Get("COOKIE_SECURE", "0") == "1", SameSite: http.SameSiteLaxMode, MaxAge: 86400})
}

func EndSession(r *http.Request) {
	c, e := r.Cookie("ticket_session")
	if e != nil {
		return
	}
	key := sha256.Sum256([]byte(c.Value))
	repositories.DeleteSession(hex.EncodeToString(key[:]))
}
