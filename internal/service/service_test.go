package service_test

import (
	"context"
	"testing"

	"social_media/internal/entity"
	"social_media/internal/repository/mocks"
	"social_media/internal/router/dto"
	"social_media/internal/service"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func setupTest() (*service.Service, *mocks.MockRepository) {
	mockRepo := new(mocks.MockRepository)
	logger := logrus.New()
	logger.SetLevel(logrus.PanicLevel)
	svc := service.New(mockRepo, logger)
	return svc, mockRepo
}

func TestService_Register(t *testing.T) {
	svc, mockRepo := setupTest()

	req := &dto.RegisterRequest{
		FirstName: "Anna",
		LastName:  "Gagikyan",
		Email:     "anna@test.com",
		Age:       30,
		Password:  "password123",
	}

	password := "password123"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockRepo.On("GetUserByEmail", mock.Anything, req.Email).Return(nil, nil).Once()

	mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
		return u.FirstName == "Anna" && u.Email == "anna@test.com" && len(u.PasswordHash) > 0
	})).Return(&entity.User{
		ID:           uuid.New(),
		FirstName:    "Anna",
		LastName:     "Gagikyan",
		Email:        "anna@test.com",
		PasswordHash: string(hash),
	}, nil)

	user, err := svc.Register(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.FirstName, user.FirstName)
	assert.NotEqual(t, req.Password, user.PasswordHash)

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestService_Register_DuplicateEmail(t *testing.T) {
	svc, mockRepo := setupTest()

	req := &dto.RegisterRequest{
		FirstName: "Anna",
		LastName:  "Gagikyan",
		Email:     "duplicate@test.com",
		Password:  "password123",
		Age:       25,
	}

	mockRepo.On("GetUserByEmail", mock.Anything, req.Email).Return(&entity.User{Email: req.Email}, nil).Once()

	user, err := svc.Register(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.ErrorIs(t, err, service.ErrEmailAlreadyExists)

	mockRepo.AssertExpectations(t)
}

func TestService_Register_ValidationErrors(t *testing.T) {
	svc, _ := setupTest()

	tests := []struct {
		name    string
		req     *dto.RegisterRequest
		wantErr error
	}{
		{"EmptyFirstName", &dto.RegisterRequest{FirstName: "", LastName: "Gagikyan", Email: "a@b.com", Password: "pass12", Age: 20}, service.ErrFirstNameEmpty},
		{"EmptyLastName", &dto.RegisterRequest{FirstName: "Anna", LastName: "", Email: "a@b.com", Password: "pass12", Age: 20}, service.ErrLastNameEmpty},
		{"InvalidEmail", &dto.RegisterRequest{FirstName: "Anna", LastName: "Gagikyan", Email: "notanemail", Password: "pass12", Age: 20}, service.ErrInvalidEmail},
		{"ShortPassword", &dto.RegisterRequest{FirstName: "Anna", LastName: "Gagikyan", Email: "a@b.com", Password: "abc", Age: 20}, service.ErrPasswordTooShort},
		{"InvalidAge", &dto.RegisterRequest{FirstName: "Anna", LastName: "Gagikyan", Email: "a@b.com", Password: "pass12", Age: -1}, service.ErrInvalidAge},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Register(context.Background(), tt.req)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestService_Login_Success(t *testing.T) {
	svc, mockRepo := setupTest()

	password := "mypassword"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	existingUser := &entity.User{
		ID:           uuid.New(),
		Email:        "test@test.com",
		PasswordHash: string(hash),
	}

	req := &dto.LoginRequest{
		Email:    "test@test.com",
		Password: "mypassword",
	}

	mockRepo.On("GetUserByEmail", mock.Anything, "test@test.com").Return(existingUser, nil)

	token, returnedUser, err := svc.Login(context.Background(), req)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, existingUser.ID, returnedUser.ID)

	mockRepo.AssertExpectations(t)
}

func TestService_Login_InvalidPassword(t *testing.T) {
	svc, mockRepo := setupTest()

	password := "correctpassword"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	existingUser := &entity.User{
		ID:           uuid.New(),
		Email:        "test@test.com",
		PasswordHash: string(hash),
	}

	req := &dto.LoginRequest{
		Email:    "test@test.com",
		Password: "wrongpassword",
	}

	mockRepo.On("GetUserByEmail", mock.Anything, "test@test.com").Return(existingUser, nil)

	token, returnedUser, err := svc.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Nil(t, returnedUser)

	mockRepo.AssertExpectations(t)
}

func TestService_SearchUsers(t *testing.T) {
	svc, mockRepo := setupTest()

	currentUserID := uuid.New()
	firstName := "Anna"
	lastName := "Gagikyan"
	age := 30

	expectedUsers := []*entity.User{
		{ID: uuid.New(), FirstName: "Anna", LastName: "Gagikyan", Age: 30},
	}

	mockRepo.On("SearchUsers", mock.Anything, firstName, lastName, age, currentUserID).Return(expectedUsers, nil).Once()

	users, err := svc.SearchUsers(context.Background(), firstName, lastName, age, currentUserID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	mockRepo.AssertExpectations(t)
}

func TestService_UpdateProfile(t *testing.T) {
	t.Run("UpdateAllFields", func(t *testing.T) {
		svc, mockRepo := setupTest()
		userID := uuid.New()

		existingUser := &entity.User{
			ID:        userID,
			FirstName: "Anna",
			LastName:  "Gagikyan",
			Age:       20,
			Email:     "anna@test.com",
		}
		newName := "Ani"
		newLast := "Anyan"
		newAge := 25
		newPass := "newpassword123"
		newEmail := "new@test.com"

		req := &dto.UpdateProfileRequest{
			FirstName: &newName,
			LastName:  &newLast,
			Age:       &newAge,
			Password:  &newPass,
			Email:     &newEmail,
		}

		mockRepo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil).Once()
		mockRepo.On("GetUserByEmail", mock.Anything, newEmail).Return(nil, nil).Once()
		mockRepo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
			return u.FirstName == newName &&
				u.LastName == newLast &&
				u.Age == newAge &&
				u.Email == newEmail &&
				u.ID == userID
		})).Return(nil).Once()

		err := svc.UpdateProfile(context.Background(), userID, req)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("PartialUpdate", func(t *testing.T) {
		svc, mockRepo := setupTest()
		userID := uuid.New()

		existingUser := &entity.User{
			ID:  userID,
			Age: 20,
		}
		newAge := 30
		req := &dto.UpdateProfileRequest{Age: &newAge}

		mockRepo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil).Once()
		mockRepo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *entity.User) bool {
			return u.Age == newAge && u.ID == userID
		})).Return(nil).Once()

		err := svc.UpdateProfile(context.Background(), userID, req)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NoChanges_SameFieldValues", func(t *testing.T) {
		svc, mockRepo := setupTest()
		userID := uuid.New()

		oldName := "Anna"
		existingUser := &entity.User{ID: userID, FirstName: oldName}
		req := &dto.UpdateProfileRequest{FirstName: &oldName}

		mockRepo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil).Once()

		err := svc.UpdateProfile(context.Background(), userID, req)
		assert.NoError(t, err)
		mockRepo.AssertNotCalled(t, "UpdateUser")
		mockRepo.AssertExpectations(t)
	})

	t.Run("SamePassword_NoChange", func(t *testing.T) {
		svc, mockRepo := setupTest()
		userID := uuid.New()

		samePass := "password123"
		hash, _ := bcrypt.GenerateFromPassword([]byte(samePass), bcrypt.DefaultCost)
		existingUser := &entity.User{ID: userID, PasswordHash: string(hash)}
		req := &dto.UpdateProfileRequest{Password: &samePass}

		mockRepo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil).Once()

		err := svc.UpdateProfile(context.Background(), userID, req)
		assert.NoError(t, err)
		mockRepo.AssertNotCalled(t, "UpdateUser")
		mockRepo.AssertExpectations(t)
	})

	t.Run("EmailDuplicateCheck", func(t *testing.T) {
		svc, mockRepo := setupTest()
		userID := uuid.New()

		existingUser := &entity.User{ID: userID, Email: "anna@test.com"}
		newEmail := "taken@test.com"
		req := &dto.UpdateProfileRequest{Email: &newEmail}

		mockRepo.On("GetUserByID", mock.Anything, userID).Return(existingUser, nil).Once()
		mockRepo.On("GetUserByEmail", mock.Anything, newEmail).Return(&entity.User{ID: uuid.New(), Email: newEmail}, nil).Once()

		err := svc.UpdateProfile(context.Background(), userID, req)
		assert.ErrorIs(t, err, service.ErrEmailTaken)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_SendFriendRequest_Success(t *testing.T) {
	svc, mockRepo := setupTest()

	senderID := uuid.New()
	receiverID := uuid.New()

	mockRepo.On("GetFriendshipByParticipants", mock.Anything, senderID, receiverID).Return(nil, nil).Once()
	mockRepo.On("CreateFriendRequest", mock.Anything, mock.MatchedBy(func(fr *entity.Friend) bool {
		return fr.SenderID == senderID && fr.ReceiverID == receiverID && fr.Status == entity.StatusPending
	})).Return(&entity.Friend{
		ID:         uuid.New(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Status:     entity.StatusPending,
	}, nil)

	err := svc.SendFriendRequest(context.Background(), senderID, receiverID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_SendFriendRequest_SameUser(t *testing.T) {
	svc, _ := setupTest()

	userID := uuid.New()
	err := svc.SendFriendRequest(context.Background(), userID, userID)

	assert.ErrorIs(t, err, service.ErrCannotSelfRequest)
}

func TestService_SendFriendRequest_Duplicate(t *testing.T) {
	svc, mockRepo := setupTest()

	senderID := uuid.New()
	receiverID := uuid.New()

	existingRequest := &entity.Friend{
		ID:         uuid.New(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Status:     entity.StatusPending,
	}

	mockRepo.On("GetFriendshipByParticipants", mock.Anything, senderID, receiverID).Return(existingRequest, nil).Once()

	err := svc.SendFriendRequest(context.Background(), senderID, receiverID)

	assert.ErrorIs(t, err, service.ErrRequestPending)
	mockRepo.AssertExpectations(t)
}

func TestService_SendFriendRequest_AfterDecline(t *testing.T) {
	svc, mockRepo := setupTest()

	senderID := uuid.New()
	receiverID := uuid.New()

	declined := &entity.Friend{
		ID:         uuid.New(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Status:     entity.StatusDeclined,
	}

	mockRepo.On("GetFriendshipByParticipants", mock.Anything, senderID, receiverID).Return(declined, nil).Once()
	mockRepo.On("RemoveFriend", mock.Anything, senderID, receiverID).Return(nil).Once()
	mockRepo.On("CreateFriendRequest", mock.Anything, mock.MatchedBy(func(fr *entity.Friend) bool {
		return fr.SenderID == senderID && fr.ReceiverID == receiverID && fr.Status == entity.StatusPending
	})).Return(&entity.Friend{
		ID:         uuid.New(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Status:     entity.StatusPending,
	}, nil).Once()

	err := svc.SendFriendRequest(context.Background(), senderID, receiverID)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_AcceptFriendRequest(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc, mockRepo := setupTest()

		receiverID := uuid.New()
		requestID := uuid.New()

		fr := &entity.Friend{
			ID:         requestID,
			SenderID:   uuid.New(),
			ReceiverID: receiverID,
			Status:     entity.StatusPending,
		}

		mockRepo.On("GetFriendRequest", mock.Anything, requestID).Return(fr, nil).Once()
		mockRepo.On("UpdateFriendRequest", mock.Anything, requestID, entity.StatusAccepted).Return(nil).Once()

		err := svc.AcceptFriendRequest(context.Background(), requestID, receiverID)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		svc, mockRepo := setupTest()

		requestID := uuid.New()
		fr := &entity.Friend{
			ID:         requestID,
			SenderID:   uuid.New(),
			ReceiverID: uuid.New(),
			Status:     entity.StatusPending,
		}

		mockRepo.On("GetFriendRequest", mock.Anything, requestID).Return(fr, nil).Once()

		err := svc.AcceptFriendRequest(context.Background(), requestID, uuid.New())
		assert.ErrorIs(t, err, service.ErrUnauthorized)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotPending", func(t *testing.T) {
		svc, mockRepo := setupTest()

		receiverID := uuid.New()
		requestID := uuid.New()

		fr := &entity.Friend{
			ID:         requestID,
			SenderID:   uuid.New(),
			ReceiverID: receiverID,
			Status:     entity.StatusAccepted,
		}

		mockRepo.On("GetFriendRequest", mock.Anything, requestID).Return(fr, nil).Once()

		err := svc.AcceptFriendRequest(context.Background(), requestID, receiverID)
		assert.ErrorIs(t, err, service.ErrNotPending)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_DeclineFriendRequest(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc, mockRepo := setupTest()

		receiverID := uuid.New()
		requestID := uuid.New()

		fr := &entity.Friend{
			ID:         requestID,
			SenderID:   uuid.New(),
			ReceiverID: receiverID,
			Status:     entity.StatusPending,
		}

		mockRepo.On("GetFriendRequest", mock.Anything, requestID).Return(fr, nil).Once()
		mockRepo.On("UpdateFriendRequest", mock.Anything, requestID, entity.StatusDeclined).Return(nil).Once()

		err := svc.DeclineFriendRequest(context.Background(), requestID, receiverID)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		svc, mockRepo := setupTest()

		requestID := uuid.New()
		fr := &entity.Friend{
			ID:         requestID,
			SenderID:   uuid.New(),
			ReceiverID: uuid.New(),
			Status:     entity.StatusPending,
		}

		mockRepo.On("GetFriendRequest", mock.Anything, requestID).Return(fr, nil).Once()

		err := svc.DeclineFriendRequest(context.Background(), requestID, uuid.New())
		assert.ErrorIs(t, err, service.ErrDeclineUnauthorized)
		mockRepo.AssertExpectations(t)
	})
}

func TestService_ListFriendRequests(t *testing.T) {
	svc, mockRepo := setupTest()

	userID := uuid.New()
	expected := []*entity.Friend{
		{ID: uuid.New(), SenderID: uuid.New(), ReceiverID: userID, Status: entity.StatusPending},
	}

	mockRepo.On("ListFriendRequests", mock.Anything, userID).Return(expected, nil).Once()

	result, err := svc.ListFriendRequests(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestService_ListFriends(t *testing.T) {
	svc, mockRepo := setupTest()

	userID := uuid.New()
	expected := []*entity.User{
		{ID: uuid.New(), FirstName: "Anna", LastName: "Gagikyan"},
		{ID: uuid.New(), FirstName: "Ani", LastName: "Anyan"},
	}

	mockRepo.On("ListFriends", mock.Anything, userID).Return(expected, nil).Once()

	result, err := svc.ListFriends(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestService_RemoveFriend(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc, mockRepo := setupTest()

		userID := uuid.New()
		friendID := uuid.New()

		mockRepo.On("RemoveFriend", mock.Anything, userID, friendID).Return(nil).Once()

		err := svc.RemoveFriend(context.Background(), userID, friendID)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		svc, mockRepo := setupTest()

		userID := uuid.New()
		friendID := uuid.New()

		mockRepo.On("RemoveFriend", mock.Anything, userID, friendID).Return(assert.AnError).Once()

		err := svc.RemoveFriend(context.Background(), userID, friendID)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
