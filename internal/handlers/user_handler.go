package handlers

import "learning/internal/services"

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	 return &UserHandler{
		userService: userService,
	 }
}

