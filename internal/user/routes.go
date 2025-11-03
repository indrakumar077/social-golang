package user

import (
	"learning/internal/database"

	"github.com/gorilla/mux"
)

// RegisterRoutes is a convenience wrapper when you already have a Handler
func RegisterRoutes(r *mux.Router, h *Handler) {
	h.RegisterRoutes(r)
}

// Register composes repository -> service -> handler and registers routes
func Register(r *mux.Router, db *database.DataBase) {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)
	h.RegisterRoutes(r)
}
