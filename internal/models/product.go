package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID 			uuid.UUID `json:"id"`
	FarmerID	uuid.UUID `json:"farmer_id"`
	Name		string	  `json:"name"`
	UnitType	string    `json:"unit_type"`
	Description	string    `json:"description"`
	Version		int       `json:"version"`
	Deleted		bool      `json:"-"`
	CreatedAt	time.Time `json:"created_at"`
	UpdatedAt	time.Time `json:"updated_at"`
}