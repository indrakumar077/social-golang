package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Name       string    `json:"name"`
	Password   string    `json:"-"`
	MiddleName *string   `json:"middle_name,omitempty"`
	Surname    *string   `json:"surname,omitempty"`
	Bio        *string   `json:"bio,omitempty"`
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Name       string    `json:"name"`
	MiddleName *string   `json:"middle_name,omitempty"`
	Surname    *string   `json:"surname,omitempty"`
	Bio        *string   `json:"bio,omitempty"`
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func ToUserReponse(u *User) *UserResponse {
	return &UserResponse{
		ID:         u.ID,
		Username:   u.Username,
		Email:      u.Email,
		Phone:      u.Phone,
		Name:       u.Name,
		MiddleName: u.MiddleName,
		Surname:    u.Surname,
		Bio:        u.Bio,
		Active:     u.Active,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}
