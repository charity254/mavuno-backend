package models

import (
	"time"

	"github.com/google/uuid"
)

type SupplyLocation struct {
	ID              uuid.UUID `json:"id"`
	FarmerID        uuid.UUID `json:"farmer_id"`
	Name            string    `json:"name"`
	ContactPerson   string    `json:"contact_person"`
	PhoneNumber     string    `json:"phone_number"`
	LocationAddress string    `json:"location_address"`
	Notes           string    `json:"notes"`
	Version         int       `json:"version"`
	Deleted         bool      `json:"-"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}