package user

import (
	"context"
	"fmt"
	"learning/internal/database"
	"time"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *database.DataBase
}

// Ensure Repository implements the expected interface
var _ RepositoryInterface = (*Repository)(nil)

// RepositoryInterface defines persistence operations for users
type RepositoryInterface interface {
	CreateUser(ctx context.Context, user *CreateUserRequest, hashedPassword string) (*User, error)
	GetUserById(ctx context.Context, id int) (*User, error)
}

// NewRepository creates a new user repository
func NewRepository(db *database.DataBase) *Repository {
	return &Repository{db: db}
}

// scanUserFromRow scans a database row into a User model
func (r *Repository) scanUserFromRow(row pgx.Row) (*User, error) {
	var user User

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

// CreateUser creates a new user in the database
func (r *Repository) CreateUser(ctx context.Context, user *CreateUserRequest, hashedPassword string) (*User, error) {
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

	createdUser, err := r.scanUserFromRow(row)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}

// GetUserById retrieves a user by ID from the database
func (r *Repository) GetUserById(ctx context.Context, id int) (*User, error) {
	query := `
        SELECT id, username, email, name, password, middle_name, surname, bio, active, created_at, updated_at
        FROM users
        WHERE id = $1 AND active = true
    `

	row := r.db.Pool.QueryRow(ctx, query, id)

	user, err := r.scanUserFromRow(row)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
