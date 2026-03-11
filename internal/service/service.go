package service

import (
	"context"

	"social_media/internal/entity"
	"social_media/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Service struct {
	repo repository.Repository
	log  *logrus.Logger
}

func New(repo repository.Repository, log *logrus.Logger) *Service {
	return &Service{repo: repo, log: log}
}

func (s *Service) Register(ctx context.Context, u *entity.User) (*entity.User, error) {
	s.log.WithFields(logrus.Fields{"email": u.Email}).Info("registering user")
	return s.repo.CreateUser(ctx, u)
}

func (s *Service) SearchUsers(ctx context.Context, firstName, lastName string, age int) ([]*entity.User, error) {
	s.log.WithFields(logrus.Fields{
		"first_name": firstName,
		"last_name":  lastName,
		"age":        age,
	}).Info("searching users")
	return s.repo.SearchUsers(ctx, firstName, lastName, age)
}

func (s *Service) SendFriendRequest(ctx context.Context, senderID, receiverID uuid.UUID) error {
	s.log.WithFields(logrus.Fields{
		"sender":   senderID,
		"receiver": receiverID,
	}).Info("sending friend request")
	fr := &entity.Friend{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Status:     entity.StatusPending,
	}
	_, err := s.repo.CreateFriendRequest(ctx, fr)
	return err
}
