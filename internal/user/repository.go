package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRespository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, u *User) (*User, error) {
	query := `INSERT INTO users (username, email, phone, name, password, middle_name, surname, bio, active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
		u.Username, u.Email, u.Phone, u.Name, u.Password,
		u.MiddleName, u.Surname, u.Bio, u.Active,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}
