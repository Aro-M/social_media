package service

import (
	"context"
	"strings"
	"time"

	"social_media/internal/config"
	"social_media/internal/entity"
	"social_media/internal/repository"
	"social_media/internal/router/dto"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo repository.Repository
	log  *logrus.Logger
}

func New(repo repository.Repository, log *logrus.Logger) *Service {
	return &Service{repo: repo, log: log}
}

// Register saves a new user with a hashed password.
func (s *Service) Register(ctx context.Context, req *dto.RegisterRequest) (*entity.User, error) {
	s.log.WithFields(logrus.Fields{"email": req.Email}).Info("registering user")

	if strings.TrimSpace(req.FirstName) == "" {
		return nil, ErrFirstNameEmpty
	}
	if strings.TrimSpace(req.LastName) == "" {
		return nil, ErrLastNameEmpty
	}
	if strings.TrimSpace(req.Email) == "" || !strings.Contains(req.Email, "@") {
		return nil, ErrInvalidEmail
	}
	if len(req.Password) < 6 {
		return nil, ErrPasswordTooShort
	}
	if req.Age < 0 || req.Age > 150 {
		return nil, ErrInvalidAge
	}

	existing, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	u := &entity.User{
		ID:           uuid.New(),
		FirstName:    strings.TrimSpace(req.FirstName),
		LastName:     strings.TrimSpace(req.LastName),
		Age:          req.Age,
		Email:        strings.TrimSpace(req.Email),
		PasswordHash: string(hashedBytes),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	return s.repo.CreateUser(ctx, u)
}

// Login checks credentials and returns a JWT token.
func (s *Service) Login(ctx context.Context, req *dto.LoginRequest) (string, *entity.User, error) {
	s.log.WithFields(logrus.Fields{"email": req.Email}).Info("login user")

	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", nil, err 
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return "", nil, err 
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(config.JWTSecret())
	if err != nil {
		return "", nil, err
	}

	return tokenString, user, nil
}

// SearchUsers filters users and excludes the session user.
func (s *Service) SearchUsers(ctx context.Context, firstName, lastName string, age int, currentUserID uuid.UUID) ([]*entity.User, error) {
	s.log.WithFields(logrus.Fields{
		"first_name":      firstName,
		"last_name":       lastName,
		"age":             age,
		"current_user_id": currentUserID,
	}).Info("searching users")
	return s.repo.SearchUsers(ctx, firstName, lastName, age, currentUserID)
}

// SendFriendRequest creates a new friendship invitation.
func (s *Service) SendFriendRequest(ctx context.Context, senderID, receiverID uuid.UUID) error {
	s.log.WithFields(logrus.Fields{
		"sender":   senderID,
		"receiver": receiverID,
	}).Info("sending friend request")

	if senderID == receiverID {
		return ErrCannotSelfRequest
	}

	// Check if a relationship already exists
	existing, err := s.repo.GetFriendshipByParticipants(ctx, senderID, receiverID)
	if err != nil {
		return err
	}
	if existing != nil {
		switch existing.Status {
		case entity.StatusAccepted:
			return ErrAlreadyFriends
		case entity.StatusPending:
			return ErrRequestPending
		case entity.StatusDeclined:
			// Allow re-sending after decline by removing the old record first
			if err := s.repo.RemoveFriend(ctx, existing.SenderID, existing.ReceiverID); err != nil {
				return err
			}
		}
	}

	now := time.Now()
	fr := &entity.Friend{
		ID:         uuid.New(),
		SenderID:   senderID,
		ReceiverID: receiverID,
		Status:     entity.StatusPending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	_, err = s.repo.CreateFriendRequest(ctx, fr)
	return err
}

// ListFriendRequests returns all pending requests for a user.
func (s *Service) ListFriendRequests(ctx context.Context, userID uuid.UUID) ([]*entity.Friend, error) {
	s.log.WithFields(logrus.Fields{
		"user_id": userID,
	}).Info("listing friend requests")
	return s.repo.ListFriendRequests(ctx, userID)
}

// AcceptFriendRequest marks a request as accepted.
func (s *Service) AcceptFriendRequest(ctx context.Context, requestID, currentUserID uuid.UUID) error {
	s.log.WithField("request_id", requestID).Info("accepting friend request")

	req, err := s.repo.GetFriendRequest(ctx, requestID)
	if err != nil {
		return err
	}

	if req.ReceiverID != currentUserID {
		return ErrUnauthorized
	}

	if req.Status != entity.StatusPending {
		return ErrNotPending
	}

	return s.repo.UpdateFriendRequest(ctx, requestID, entity.StatusAccepted)
}

// DeclineFriendRequest marks a request as declined.
func (s *Service) DeclineFriendRequest(ctx context.Context, requestID, currentUserID uuid.UUID) error {
	s.log.WithField("request_id", requestID).Info("declining friend request")

	req, err := s.repo.GetFriendRequest(ctx, requestID)
	if err != nil {
		return err
	}

	if req.ReceiverID != currentUserID {
		return ErrDeclineUnauthorized
	}

	if req.Status != entity.StatusPending {
		return ErrNotPending
	}

	return s.repo.UpdateFriendRequest(ctx, requestID, entity.StatusDeclined)
}

// ListFriends gets all users in the friend list.
func (s *Service) ListFriends(ctx context.Context, userID uuid.UUID) ([]*entity.User, error) {
	s.log.WithFields(logrus.Fields{
		"user_id": userID,
	}).Info("listing friends")
	return s.repo.ListFriends(ctx, userID)
}

// RemoveFriend deletes a friendship record.
func (s *Service) RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	s.log.WithFields(logrus.Fields{
		"user_id":   userID,
		"friend_id": friendID,
	}).Info("removing friend")
	return s.repo.RemoveFriend(ctx, userID, friendID)
}

// UpdateProfile modifies user information.
func (s *Service) UpdateProfile(ctx context.Context, userID uuid.UUID, req *dto.UpdateProfileRequest) error {
	s.log.WithField("user_id", userID).Info("updating user profile")

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	changed := false

	// Helper for updating string fields
	updateStr := func(field *string, newVal *string) {
		if newVal != nil && *newVal != *field {
			*field = *newVal
			changed = true
		}
	}

	updateStr(&user.FirstName, req.FirstName)
	updateStr(&user.LastName, req.LastName)

	if req.Age != nil && *req.Age != user.Age {
		user.Age = *req.Age
		changed = true
	}

	if req.Email != nil && *req.Email != user.Email {
		if strings.TrimSpace(*req.Email) == "" || !strings.Contains(*req.Email, "@") {
			return ErrInvalidEmail
		}
		existing, err := s.repo.GetUserByEmail(ctx, *req.Email)
		if err == nil && existing != nil && existing.ID != userID {
			return ErrEmailTaken
		}
		user.Email = strings.TrimSpace(*req.Email)
		changed = true
	}

	if req.Password != nil {
		if len(*req.Password) < 6 {
			return ErrPasswordTooShort
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(*req.Password)); err != nil {
			hashedBytes, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			user.PasswordHash = string(hashedBytes)
			changed = true
		}
	}

	if !changed {
		return nil
	}

	user.UpdatedAt = time.Now()
	return s.repo.UpdateUser(ctx, user)
}
