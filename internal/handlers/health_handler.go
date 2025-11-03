package handlers

import (
	"net/http"

	"learning/internal/database"
	"learning/internal/utils"

	"github.com/gorilla/mux"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db *database.DataBase
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *database.DataBase) *HealthHandler {
	return &HealthHandler{db: db}
}

// RegisterRoutes registers health check routes
func (h *HealthHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/health", h.Health).Methods(http.MethodGet)
	r.HandleFunc("/health/ready", h.Readiness).Methods(http.MethodGet)
}

// RegisterHealth is a convenience function that composes and registers routes
func RegisterHealth(r *mux.Router, db *database.DataBase) {
	h := NewHealthHandler(db)
	h.RegisterRoutes(r)
}

// Health returns basic health status
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	utils.WriteSuccess(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

// Readiness checks if the application is ready to serve traffic
func (h *HealthHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	// Check database connection
	ctx := r.Context()
	if err := h.db.Pool.Ping(ctx); err != nil {
		utils.WriteError(w, http.StatusServiceUnavailable, "database connection failed")
		return
	}

	utils.WriteSuccess(w, http.StatusOK, map[string]string{
		"status": "ready",
	})
}
