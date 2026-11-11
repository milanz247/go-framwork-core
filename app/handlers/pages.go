package handlers

import (
	"crypto/hmac"
	"encoding/json"
	"net/http"

	"ticket-sync-mvc/app/auth"
	"ticket-sync-mvc/app/views"
	"ticket-sync-mvc/framework/env"
	"ticket-sync-mvc/framework/support"
)

func redirect(w http.ResponseWriter, r *http.Request, path string) {
	if r.Header.Get("HX-Request") == "true" {
		if path == "/login" || r.URL.Path == "/login" || r.URL.Path == "/register" {
			w.Header().Set("HX-Redirect", path)
		} else {
			b, _ := json.Marshal(map[string]string{"path": path, "target": "#workspace", "swap": "outerHTML"})
			w.Header().Set("HX-Location", string(b))
		}
		w.WriteHeader(200)
	} else {
		http.Redirect(w, r, path, http.StatusSeeOther)
	}
}
func HTML(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" && r.URL.Path != "/dashboard" && r.URL.Path != "/login" && r.URL.Path != "/register" && r.URL.Path != "/logout" {
		http.NotFound(w, r)
		return
	}
	if r.Method != "GET" && r.Method != "POST" {
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Method not allowed", 405)
		return
	}
	user := auth.SessionUser(r)
	if r.URL.Path == "/login" || r.URL.Path == "/register" {
		authPage(w, r, user)
		return
	}
	if user == nil {
		redirect(w, r, "/login")
		return
	}
	if r.URL.Path == "/logout" {
		if r.Method != "POST" {
			w.Header().Set("Allow", "POST")
			http.Error(w, "POST required", 405)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form", 400)
			return
		}
		if !hmac.Equal([]byte(r.FormValue("csrf")), []byte(support.Str(user["csrf"]))) {
			http.Error(w, "Refresh the page and try again", 403)
			return
		}
		handleLogout(w, r)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "GET required", 405)
		return
	}
	if r.URL.Path == "/" {
		redirect(w, r, "/dashboard")
		return
	}
	views.Render(w, r, 200, support.M{"User": user, "CSRF": user["csrf"], "Page": "dashboard", "Title": "Dashboard", "Dev": env.Get("APP_ENV", "production") == "development"})
}
func authPage(w http.ResponseWriter, r *http.Request, u support.M) {
	if u != nil {
		redirect(w, r, "/dashboard")
		return
	}
	token := ""
	if c, e := r.Cookie("ticket_login_csrf"); e == nil {
		token = c.Value
	}
	if len(token) != 64 {
		token = support.Random()
		http.SetCookie(w, &http.Cookie{Name: "ticket_login_csrf", Value: token, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode, Secure: env.Get("COOKIE_SECURE", "0") == "1", MaxAge: 3600})
	}
	register := r.URL.Path == "/register"
	data := support.M{"Auth": true, "Register": register, "CSRF": token, "Title": "Sign in", "Dev": env.Get("APP_ENV", "production") == "development"}
	code := 200
	if register {
		data["Title"] = "Create account"
	}
	if r.Method == "POST" {
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		if e := r.ParseForm(); e != nil {
			http.Error(w, "Invalid form", 400)
			return
		}
		if !hmac.Equal([]byte(token), []byte(r.FormValue("csrf"))) {
			http.Error(w, "Refresh the page and try again", 403)
			return
		}
		v := support.M{"email": r.FormValue("email"), "password": r.FormValue("password")}
		id, e := auth.Authenticate(v, register)
		if e == nil {
			auth.StartSession(w, id)
			redirect(w, r, "/dashboard")
			return
		}
		if p, ok := e.(*auth.Problem); ok {
			code = p.Status
			data["Error"] = p.Message
		} else {
			panic(e)
		}
		data["Email"] = v["email"]
	}
	views.Render(w, r, code, data)
}
