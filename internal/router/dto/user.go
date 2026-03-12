package dto

import (
	"github.com/google/uuid"
)

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Age       int       `json:"age"`
	Email     string    `json:"email"`
}

type SearchRequest struct {
	FirstName string `query:"first_name"`
	LastName  string `query:"last_name"`
	Age       int    `query:"age"`
}

type UpdateProfileRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Age       *int    `json:"age"`
	Email     *string `json:"email"`
	Password  *string `json:"password"`
}
