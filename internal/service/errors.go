package service

import "errors"

var (
	ErrEmailTaken          = errors.New("email already taken")
	ErrEmailAlreadyExists  = errors.New("user with this email already exists")
	ErrAlreadyFriends      = errors.New("already friends")
	ErrRequestPending      = errors.New("friend request already pending")
	ErrCannotSelfRequest   = errors.New("cannot send friend request to yourself")
	ErrUnauthorized        = errors.New("unauthorized: only the receiver can accept the request")
	ErrDeclineUnauthorized = errors.New("unauthorized: only the receiver can decline the request")
	ErrNotPending          = errors.New("request is not pending")

	ErrInvalidEmail     = errors.New("email is required and must be valid")
	ErrInvalidAge       = errors.New("age must be between 0 and 150")
	ErrPasswordTooShort = errors.New("password must be at least 6 characters")
	ErrFirstNameEmpty   = errors.New("first name is required")
	ErrLastNameEmpty    = errors.New("last name is required")
)
