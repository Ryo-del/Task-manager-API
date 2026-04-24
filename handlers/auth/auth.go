package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	authservice "taskmanager/service/authService"
)

type Server struct {
	service authservice.AuthServicer
}

func NewServer(service authservice.AuthServicer) *Server {
	return &Server{service: service}
}

type credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var creds credentials
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}
	} else {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Invalid form body", http.StatusBadRequest)
			return
		}
		creds.Login = r.FormValue("login")
		creds.Password = r.FormValue("password")
	}

	if creds.Login == "" || creds.Password == "" {
		http.Error(w, "Login and password are required", http.StatusBadRequest)
		return
	}

	token, err := s.service.Login(r.Context(), creds.Login, creds.Password)
	if err != nil {
		slog.Error("login failed", "err", err)
		http.Error(w, "invalid login or password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})

}

func (s *Server) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if creds.Login == "" || creds.Password == "" {
		http.Error(w, "login and password are required", http.StatusBadRequest)
		return
	}
	err := s.service.Register(r.Context(), creds.Login, creds.Password)
	if err != nil {
		slog.Error("reg failed", "err", err)
		http.Error(w, "could not register", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "user registered successfully"})
}
