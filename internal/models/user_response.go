package models

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
