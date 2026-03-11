package dto

import (
	"github.com/google/uuid"
)

type RegisterRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Age       int    `json:"age" validate:"required,min=1"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
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
