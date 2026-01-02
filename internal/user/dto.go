package user

type CreateUserRequest struct {
	Username   string  `json:"username" validate:"required,min=3,max=50"`
	Email      string  `json:"email" validate:"required,email"`
	Phone      string  `json:"phone" validate:"required,min=10,max=15"`
	Name       string  `json:"name" validate:"required,min=2,max=100"`
	Password   string  `json:"password" validate:"required,min=6"`
	MiddleName *string `json:"middle_name,omitempty" validate:"omitempty,max=100"`
	Surname    *string `json:"surname,omitempty" validate:"omitempty,max=100"`
	Bio        *string `json:"bio,omitempty"`
}
