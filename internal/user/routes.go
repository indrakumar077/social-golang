package user

import (
	"learning/internal/storage"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(router *mux.Router, pool *pgxpool.Pool, s3Client *storage.S3Client) {

	repo := NewRespository(pool)
	service := NewService(repo, s3Client)
	handler := NewHandler(service)

	userRouter := router.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("", handler.Create).Methods("POST")
}
