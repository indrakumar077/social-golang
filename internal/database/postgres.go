package database

import (
	"context"
	"fmt"
	"learning/internal/config"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DataBase struct {
	Pool *pgxpool.Pool
}


func New(cfg *config.Config) (*DataBase, error) {
	// Validate database configuration
	if err := cfg.DataBase.Validate(); err != nil {
		return nil, fmt.Errorf("invalid database configuration: %w", err)
	}

	poolConfig, err := pgxpool.ParseConfig(cfg.DataBase.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database configuration: %w", err)
	}

	poolConfig.MaxConns = cfg.DataBase.MaxConn
	poolConfig.MinConns = cfg.DataBase.MinConn
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = time.Minute * 30

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("database connection test failed: %w", err)
	}

	log.Println("Database connected successfully")

	return &DataBase{Pool: pool}, nil
}

func (db *DataBase) Close() {
	db.Pool.Close()
	log.Println("Database connection closed")
}