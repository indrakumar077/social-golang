package user

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router *mux.Router, pool *pgxpool.Pool) {

	repo := NewRespository(pool)
	service := NewService(repo)
	handler := NewHandler(service)

	router.PathPrefix("/users").Subrouter()
	router.HandleFunc("", handler.Create).Methods("POST")
}
