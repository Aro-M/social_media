package service

import (
	"context"
	"social_media/internal/entity"
	"social_media/internal/router/dto"

	"github.com/google/uuid"
)

type ServiceInterface interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*entity.User, error)
	Login(ctx context.Context, req *dto.LoginRequest) (string, *entity.User, error)
	SearchUsers(ctx context.Context, firstName, lastName string, age int, currentUserID uuid.UUID) ([]*entity.User, error)
	SendFriendRequest(ctx context.Context, senderID, receiverID uuid.UUID) error
	ListFriendRequests(ctx context.Context, userID uuid.UUID) ([]*entity.Friend, error)
	AcceptFriendRequest(ctx context.Context, requestID, currentUserID uuid.UUID) error
	DeclineFriendRequest(ctx context.Context, requestID, currentUserID uuid.UUID) error
	ListFriends(ctx context.Context, userID uuid.UUID) ([]*entity.User, error)
	RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) error
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateProfileRequest) error
}
