package handlers

import (
	"net/http"

	"ticket-sync-mvc/app/auth"
	"ticket-sync-mvc/framework/http/response"
	"ticket-sync-mvc/framework/support"
)

func handleAuth(w http.ResponseWriter, r *http.Request, path string) {
	if r.Method != "POST" {
		response.Fail(w, 405, "POST required")
		return
	}
	v, e := response.Decode(w, r)
	if e != nil {
		response.Fail(w, 400, "Invalid JSON")
		return
	}
	id, e := auth.Authenticate(v, path == "register")
	if e != nil {
		if p, ok := e.(*auth.Problem); ok {
			response.Fail(w, p.Status, p.Message)
			return
		}
		panic(e)
	}
	auth.StartSession(w, id)
	response.JSON(w, support.M{"ok": true})
}
