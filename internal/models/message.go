package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID `json:"id"`
	ListingID uuid.UUID `json:"listing_id"`
	SenderID  uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	Deleted   bool      `json:"-"`
}