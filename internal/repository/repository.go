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
	SearchUsers(ctx context.Context, firstName, lastName string, age int, currentUserID uuid.UUID) ([]*entity.User, error)
	UpdateUser(ctx context.Context, u *entity.User) error

	CreateFriendRequest(ctx context.Context, fr *entity.Friend) (*entity.Friend, error)
	GetFriendRequest(ctx context.Context, id uuid.UUID) (*entity.Friend, error)
	UpdateFriendRequest(ctx context.Context, id uuid.UUID, status entity.FriendStatus) error
	ListFriendRequests(ctx context.Context, userID uuid.UUID) ([]*entity.Friend, error)
	ListFriends(ctx context.Context, userID uuid.UUID) ([]*entity.User, error)
	RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) error
	GetFriendshipByParticipants(ctx context.Context, userID1, userID2 uuid.UUID) (*entity.Friend, error)

	Close() error
}
