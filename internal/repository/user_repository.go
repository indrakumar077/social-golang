package repository

import (
	"context"
	"fmt"
	"learning/internal/database"
	"learning/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *database.DataBase
}

func NewUserRepository(db *database.DataBase) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) scanUserFromRow(row pgx.Row) (*models.User, error) {
	var user models.User
	
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Name,
		&user.Password,
		&user.MiddleName,
		&user.Surname,
		&user.Bio,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to scan user: %w", err)
	}

	return &user, nil
}


func (r *UserRepository) CreateUser(ctx context.Context, user *models.CreateUserRequest, hashedPassword string) (*models.User,error){
	query := `
		INSERT INTO users (username, email, name, password, middle_name, surname, bio, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, true, $8, $9)
		RETURNING id, username, email, name, password, middle_name, surname, bio, active, created_at, updated_at
	`

	now := time.Now()
	row := r.db.Pool.QueryRow(ctx, query,  
		user.Username,
		user.Email,
		user.Name,
		hashedPassword,
		user.MiddleName,
		user.Surname,
		user.Bio,
		now,
		now,
	)

	createdUser , err := r.scanUserFromRow(row)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}


func (r *UserRepository) GetUserById(ctx context.Context, id int) (*models.User,error){
	 	query := `
		SELECT id, username, email, name, password, middle_name, surname, bio, active, created_at, updated_at
		FROM users
		WHERE id = $1 AND active = true
	`

	row  := r.db.Pool.QueryRow(ctx,query,id);

	user , err := r.scanUserFromRow(row)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}