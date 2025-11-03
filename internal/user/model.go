package user

import "time"

// User represents a user entity in the system
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

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Username   string  `json:"username" validate:"required,min=3,max=50"`
	Email      string  `json:"email" validate:"required,email"`
	Name       string  `json:"name" validate:"required"`
	Password   string  `json:"password" validate:"required,min=6"`
	MiddleName *string `json:"middle_name,omitempty"`
	Surname    *string `json:"surname,omitempty"`
	Bio        *string `json:"bio,omitempty"`
}

// UserResponse represents the user data returned in API responses (password excluded)
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

// ToUserResponse converts a User to UserResponse (excluding password)
func ToUserResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		Name:       user.Name,
		MiddleName: user.MiddleName,
		Surname:    user.Surname,
		Bio:        user.Bio,
		Active:     user.Active,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}
