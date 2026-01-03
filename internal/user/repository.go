package user

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

// scanUser scans a pgx.Row into a User struct
func scanUser(row pgx.Row) (*User, error) {
	var user User
	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.Phone, &user.Name,
		&user.Password, &user.MiddleName, &user.Surname, &user.Bio,
		&user.Active, &user.CreatedAt, &user.UpdatedAt, &user.ProfilePhotoURL,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `SELECT id, username, email, phone, name, password, middle_name, surname, bio, 
		active, created_at, updated_at, profile_photo_url
		FROM users WHERE id = $1`

	user, err := scanUser(r.pool.QueryRow(ctx, query, id))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, username, email, phone, name, password, middle_name, surname, bio, 
		active, created_at, updated_at, profile_photo_url
		FROM users WHERE email = $1`

	user, err := scanUser(r.pool.QueryRow(ctx, query, email))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (*User, error) {
	query := `SELECT id, username, email, phone, name, password, middle_name, surname, bio, 
		active, created_at, updated_at, profile_photo_url
		FROM users WHERE username = $1`

	user, err := scanUser(r.pool.QueryRow(ctx, query, username))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}
