package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"remoteChange/internal/domain/config"
	"remoteChange/internal/middleware"
	"remoteChange/internal/model"
	"strconv"
)

type CfgHandler struct {
	service *config.Service
}

func NewCfgHandler(service *config.Service) *CfgHandler {
	return &CfgHandler{service: service}
}

func (h *CfgHandler) RegisterRoutes(r *mux.Router) {
	r.Handle("/configs", middleware.AuthMiddleware(http.HandlerFunc(h.CreateConfig))).Methods("POST")
	r.Handle("/configs/{id}", middleware.AuthMiddleware(http.HandlerFunc(h.EditConfig))).Methods("PUT")
	r.Handle("/configs/team/{teamId}", middleware.AuthMiddleware(http.HandlerFunc(h.GetConfigByTeam))).Methods("GET")
	r.Handle("/configs/{id}/rollback", middleware.AuthMiddleware(http.HandlerFunc(h.Rollback))).Methods("POST")
	r.Handle("/configs/{id}", middleware.AuthMiddleware(http.HandlerFunc(h.GetConfigById))).Methods("GET")
}

func (h *CfgHandler) CreateConfig(w http.ResponseWriter, r *http.Request) {
	var dto model.CreateConfigDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := h.service.CreateConfig(r.Context(), dto); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondWithJSON(w, http.StatusCreated, map[string]string{"message": "config created"})
}

func (h *CfgHandler) EditConfig(w http.ResponseWriter, r *http.Request) {
	var dto model.UpdateConfigDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := h.service.EditConfig(r.Context(), dto); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "config updated"})
}

func (h *CfgHandler) GetConfigByTeam(w http.ResponseWriter, r *http.Request) {
	teamId, err := strconv.ParseInt(mux.Vars(r)["teamId"], 10, 64)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid team ID")
		return
	}
	configs, err := h.service.GetConfigByTeam(teamId)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondWithJSON(w, http.StatusOK, configs)
}

func (h *CfgHandler) Rollback(w http.ResponseWriter, r *http.Request) {
	configId, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid config ID")
		return
	}
	if err := h.service.Rollback(r.Context(), configId); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondWithJSON(w, http.StatusOK, map[string]string{"message": "config rolled back"})
}

func (h *CfgHandler) GetConfigById(w http.ResponseWriter, r *http.Request) {
	configId, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "invalid config ID")
		return
	}
	config, err := h.service.GetConfigById(configId)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.respondWithJSON(w, http.StatusOK, config)
}

func (h *CfgHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *CfgHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
