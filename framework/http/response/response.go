package response

import (
	"encoding/json"
	"net/http"

	"ticket-sync-mvc/framework/support"
)

func JSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func Fail(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(support.M{"error": msg})
}

func Decode(w http.ResponseWriter, r *http.Request) (support.M, error) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	v := support.M{}
	e := json.NewDecoder(r.Body).Decode(&v)
	return v, e
}
