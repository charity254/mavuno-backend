package models

import (
	"time"

	"github.com/google/uuid"
)

// SupplyAgreement represents a recurring supply commitment a farmer makes.
type SupplyAgreement struct {
	ID               uuid.UUID `json:"id"`
	FarmerID         uuid.UUID `json:"farmer_id"`
	ProductID        uuid.UUID `json:"product_id"`
	SupplyLocationID uuid.UUID `json:"supply_location_id"`
	QtyPerDelivery   int       `json:"quantity_per_delivery"` // how much to deliver each time
	PricePerUnit     int       `json:"price_per_unit"`        // in KES cents
	DeliveryDays     []string  `json:"delivery_days"`         // e.g. ["Monday", "Wednesday"]
	DeliveryNotes    string    `json:"delivery_notes"`
	Active           bool      `json:"active"`                // can be toggled without deleting
	Version          int       `json:"version"`
	Deleted          bool      `json:"-"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}