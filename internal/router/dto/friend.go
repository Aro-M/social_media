package dto

import (
	"github.com/google/uuid"
)

type FriendRequestDTO struct {
	ReceiverID uuid.UUID `json:"receiver_id" validate:"required"`
}

type FriendRequestResponse struct {
	ID         uuid.UUID `json:"id"`
	SenderID   uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Status     string    `json:"status"`
}

type UpdateFriendRequest struct {
	Status string `json:"status" validate:"required,oneof=accepted declined"`
}
