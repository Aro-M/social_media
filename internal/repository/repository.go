package repository

import (
	"context"

	"social_media/internal/entity"

	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(ctx context.Context, u *entity.User) (*entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	SearchUsers(ctx context.Context, firstName, lastName string, age int) ([]*entity.User, error)

	CreateFriendRequest(ctx context.Context, fr *entity.Friend) (*entity.Friend, error)
	GetFriendRequest(ctx context.Context, id uuid.UUID) (*entity.Friend, error)
	UpdateFriendRequest(ctx context.Context, id uuid.UUID, status entity.FriendStatus) error
	ListFriendRequests(ctx context.Context, userID uuid.UUID) ([]*entity.Friend, error)

	Close() error
}
