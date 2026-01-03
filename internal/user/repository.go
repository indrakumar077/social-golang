package user

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRespository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func isUniqueConstraintError(err error, field string) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// PostgreSQL unique constraint violation code
		if pgErr.Code == "23505" {
			// Check if error message contains the field name
			return strings.Contains(pgErr.ConstraintName, field) ||
				strings.Contains(pgErr.Message, field)
		}
	}
	return false
}

func (r *Repository) Create(ctx context.Context, u *User) (*User, error) {
	query := `INSERT INTO users (username, email, phone, name, password, middle_name, surname, bio, profile_photo_url, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
		u.Username, u.Email, u.Phone, u.Name, u.Password,
		u.MiddleName, u.Surname, u.Bio, u.ProfilePhotoURL, u.Active,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		// Check for unique constraint violations
		if isUniqueConstraintError(err, "username") {
			return nil, errors.New("username already exists")
		}
		if isUniqueConstraintError(err, "email") {
			return nil, errors.New("email already exists")
		}
		if isUniqueConstraintError(err, "phone") {
			return nil, errors.New("phone already exists")
		}
		return nil, err
	}
	return u, nil
}
