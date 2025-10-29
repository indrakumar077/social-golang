package main

import (
	"context"
	"learning/internal/config"
	"learning/internal/database"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting application...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}
	log.Println("Configuration loaded successfully")

	db, err := database.New(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected successfully")

	defer db.Close()

	// Initialize repository (concrete implementation)
	router := mux.NewRouter()



	server := &http.Server{
		Addr: ":" + cfg.ServerPort,
		Handler: router,
		ReadTimeout: 15*time.Second,
		WriteTimeout: 15*time.Second,
		IdleTimeout: 60*time.Second,
	}

	go func() {
		log.Println("Server starting on port:", cfg.ServerPort)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start:", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server stopped gracefully")
}