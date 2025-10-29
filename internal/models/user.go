package models

import "time"

type User struct {
	ID         int       `json:"id" db:"id"`
	Username   string    `json:"username" db:"username"`
	Email      string    `json:"email" db:"email"`
	Name       string    `json:"name" db:"name"`
	Password   string    `json:"-" db:"password"`
	MiddleName *string   `json:"middle_name,omitempty" db:"middle_name"`
	Surname    *string   `json:"surname,omitempty" db:"surname"`
	Bio        *string   `json:"bio,omitempty" db:"bio"`
	Active     bool      `json:"active" db:"active"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type CreateUserRequest struct {
	Username   string  `json:"username" validate:"required,min=3,max=50"`
	Email      string  `json:"email" validate:"required,email"`
	Name       string  `json:"name" validate:"required"`
	Password   string  `json:"password" validate:"required,min=6"`
	MiddleName *string `json:"middle_name,omitempty"`
	Surname    *string `json:"surname,omitempty"`
	Bio        *string `json:"bio,omitempty"`
}

type UserResponse struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	MiddleName *string   `json:"middle_name,omitempty"`
	Surname    *string   `json:"surname,omitempty"`
	Bio        *string   `json:"bio,omitempty"`
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}