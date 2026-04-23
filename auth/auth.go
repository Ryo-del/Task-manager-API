package auth

import (
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	login := r.FormValue("login")
	password := r.FormValue("password")

	if login == "" || password == "" {
		http.Error(w, "Login and password are required", http.StatusBadRequest)
		return
	}

}
