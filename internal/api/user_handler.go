package api

import (
	"encoding/json"
	"net/http"
	"remoteChange/internal/domain/user"

	"remoteChange/internal/model"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	service user.Service
}

func NewUserHandler(service user.Service) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user model.UserDTORegister
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := h.service.Register(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user model.UserDTOLogin
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *UserHandler) RegisterUserRoutes(router *mux.Router) {
	router.HandleFunc("/register", h.Register).Methods("POST")
	router.HandleFunc("/login", h.Login).Methods("POST")
}
