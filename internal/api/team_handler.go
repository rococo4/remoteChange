package api

import (
	"context"
	"encoding/json"
	"net/http"
	"remoteChange/internal/domain/team"
	"remoteChange/internal/middleware"
	"strconv"

	"remoteChange/internal/model"

	"github.com/gorilla/mux"
)

type TeamHandler struct {
	service team.Service
}

func NewTeamHandler(service team.Service) *TeamHandler {
	return &TeamHandler{service: service}
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var team model.TeamCreateDto
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateTeam(context.Background(), team); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *TeamHandler) EditUserInTeam(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Username string `json:"username"`
		TeamID   *int64 `json:"team_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.EditUserInTeam(request.Username, request.TeamID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TeamHandler) GetTeamById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	team, err := h.service.GetTeamById(teamID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(team)
}

func (h *TeamHandler) EditTeam(w http.ResponseWriter, r *http.Request) {
	var team model.UpdateTeamDTO
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.EditTeam(team); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TeamHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteTeam(teamID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TeamHandler) GetUserTeam(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	team, err := h.service.GetTeamForUsername(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(team)
}

func (h *TeamHandler) RegisterTeamRoutes(router *mux.Router) {
	router.Handle("/teams", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(h.CreateTeam)))).Methods("POST")
	router.Handle("/teams/{id}", middleware.AuthMiddleware(http.HandlerFunc(h.GetTeamById))).Methods("GET")
	router.Handle("/teams", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(h.EditTeam)))).Methods("PUT")
	router.Handle("/teams/{id}", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(h.DeleteTeam)))).Methods("DELETE")
	router.Handle("/teams/user", middleware.AuthMiddleware(middleware.AdminMiddleware(http.HandlerFunc(h.EditUserInTeam)))).Methods("PATCH")
	router.Handle("/teams/user/team", middleware.AuthMiddleware(http.HandlerFunc(h.GetUserTeam))).Methods("GET")
}
