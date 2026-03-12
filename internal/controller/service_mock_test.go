package controller

import (
	"context"
	"social_media/internal/entity"
	"social_media/internal/router/dto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Register(ctx context.Context, req *dto.RegisterRequest) (*entity.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockService) Login(ctx context.Context, req *dto.LoginRequest) (string, *entity.User, error) {
	args := m.Called(ctx, req)
	if args.Get(1) == nil {
		return args.String(0), nil, args.Error(2)
	}
	return args.String(0), args.Get(1).(*entity.User), args.Error(2)
}

func (m *MockService) SearchUsers(ctx context.Context, firstName, lastName string, age int, currentUserID uuid.UUID) ([]*entity.User, error) {
	args := m.Called(ctx, firstName, lastName, age, currentUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockService) SendFriendRequest(ctx context.Context, senderID, receiverID uuid.UUID) error {
	args := m.Called(ctx, senderID, receiverID)
	return args.Error(0)
}

func (m *MockService) ListFriendRequests(ctx context.Context, userID uuid.UUID) ([]*entity.Friend, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Friend), args.Error(1)
}

func (m *MockService) AcceptFriendRequest(ctx context.Context, requestID, currentUserID uuid.UUID) error {
	args := m.Called(ctx, requestID, currentUserID)
	return args.Error(0)
}

func (m *MockService) DeclineFriendRequest(ctx context.Context, requestID, currentUserID uuid.UUID) error {
	args := m.Called(ctx, requestID, currentUserID)
	return args.Error(0)
}

func (m *MockService) ListFriends(ctx context.Context, userID uuid.UUID) ([]*entity.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockService) RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	args := m.Called(ctx, userID, friendID)
	return args.Error(0)
}

func (m *MockService) UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateProfileRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}
