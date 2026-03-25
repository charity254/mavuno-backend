package models

import (
	"time"
	
	"github.com/google/uuid"
)

type ProduceEntry struct {
	ID 					uuid.UUID 	`json:"id"`
	FarmerID 			uuid.UUID 	`json:"farmer_id"`
	ProductID 			uuid.UUID 	`json:"product_id"`
	EntryDate 			time.Time 	`json:"entry_date"`
	OpeningStock 		int 		`json:"opening_stock"`  //stock at start of day
	AddedStock 			int 		`json:"added_stock"` //stock added during day   
	SoldQuantity 		int 		`json:"sold_quantity"` //stock sold
	RejectedQuantity 	int 		`json:"rejected_quantity"` //rehected or spoiled
	PricePerUnit 		int 		`json:"price_per_unit"` //price in KES cents *Frontend to convert cents into KES for display
	Notes 				string 		`json:"notes"`
	Version 			int 		`json:"version"`
	Deleted 			bool 		`json:"-"`
	CreatedAt 			time.Time 	`json:"created_at"`
	UpdatedAt 			time.Time 	`json:"updated_at"`

	//computed fields not stored in database. Calculated in the service layer and included in API responses
	TotalAvailable 		int `json:"total_available"`
	RemainingStock 		int `json:"remaining_stock"`
	RevenueGenerated 	int `json:"revenue-generated"`
}