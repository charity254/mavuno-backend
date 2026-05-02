package models

import (
	"time"

	"github.com/google/uuid"
)

type Listing struct {
	ID        uuid.UUID `json:"id"`
	FarmerID  uuid.UUID `json:"farmer_id"`
	Title     string    `json:"title"`
	Description string  `json:"description"`
	Price     int       `json:"price"`
	Quantity  int       `json:"quantity"`
	UnitType  string    `json:"unit_type"`
	Location  string    `json:"location"`
	Version   int       `json:"version"`
	Deleted   bool      `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}