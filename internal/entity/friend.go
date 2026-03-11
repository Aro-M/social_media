package entity

import (
	"time"

	"github.com/google/uuid"
)

type FriendStatus string

const (
	StatusPending  FriendStatus = "pending"
	StatusAccepted FriendStatus = "accepted"
	StatusDeclined FriendStatus = "declined"
)

type Friend struct {
	ID         uuid.UUID    `json:"id"`
	SenderID   uuid.UUID    `json:"sender_id"`
	ReceiverID uuid.UUID    `json:"receiver_id"`
	Status     FriendStatus `json:"status"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}
