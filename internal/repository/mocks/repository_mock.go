package mocks

import (
	"context"

	"social_media/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateUser(ctx context.Context, u *entity.User) (*entity.User, error) {
	args := m.Called(ctx, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockRepository) SearchUsers(ctx context.Context, firstName, lastName string, age int, currentUserID uuid.UUID) ([]*entity.User, error) {
	args := m.Called(ctx, firstName, lastName, age, currentUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockRepository) UpdateUser(ctx context.Context, u *entity.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockRepository) CreateFriendRequest(ctx context.Context, fr *entity.Friend) (*entity.Friend, error) {
	args := m.Called(ctx, fr)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Friend), args.Error(1)
}

func (m *MockRepository) GetFriendRequest(ctx context.Context, id uuid.UUID) (*entity.Friend, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Friend), args.Error(1)
}

func (m *MockRepository) UpdateFriendRequest(ctx context.Context, id uuid.UUID, status entity.FriendStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockRepository) ListFriendRequests(ctx context.Context, userID uuid.UUID) ([]*entity.Friend, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Friend), args.Error(1)
}

func (m *MockRepository) ListFriends(ctx context.Context, userID uuid.UUID) ([]*entity.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockRepository) RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	args := m.Called(ctx, userID, friendID)
	return args.Error(0)
}

func (m *MockRepository) GetFriendshipByParticipants(ctx context.Context, userID1, userID2 uuid.UUID) (*entity.Friend, error) {
	args := m.Called(ctx, userID1, userID2)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Friend), args.Error(1)
}

func (m *MockRepository) Close() error {
	args := m.Called()
	return args.Error(0)
}
