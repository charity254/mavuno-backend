package services

import (
	"fmt"
	"time"
	

	"github.com/google/uuid"
	"github.com/mavuno/mavuno-backend/internal/models"
	"github.com/mavuno/mavuno-backend/internal/storage"
)

type ProduceEntryService struct {  //handles all business logic
	eRepo 	*storage.ProduceEntryRepository
	pRepo *storage.ProductRepository
}

func NewProduceEntryService(eRepo *storage.ProduceEntryRepository, pRepo *storage.ProductRepository) *ProduceEntryService {
	return &ProduceEntryService{
		eRepo: eRepo,
		pRepo: pRepo,
	}
} 

func compute(ent *models.ProduceEntry) {  //compute calculates derived fields for an entry.Not stored in the database 
	ent.TotalAvailable = ent.OpeningStock + ent.AddedStock
	ent.RemainingStock = ent.TotalAvailable - ent.SoldQuantity - ent.RejectedQuantity
	ent.RevenueGenerated = ent.SoldQuantity * ent.PricePerUnit
}

func validate(ent *models.ProduceEntry) error {  // enforces business rules for a produce entry.
	if ent.OpeningStock < 0 {
		return fmt.Errorf("opening stock cannot be negative")
	}
	if ent.AddedStock < 0 {
		return fmt.Errorf("added stock cannot be negative")
	}
	if ent.SoldQuantity < 0 {
		return fmt.Errorf("sold quantity cannot be negative")
	}
	if ent.RejectedQuantity < 0 {
		return fmt.Errorf("rejected quantity cannot be negative")
	}
	if ent.PricePerUnit < 0 {
		return fmt.Errorf("price per unit cannot be negative")
	}

	totalAvail := ent.OpeningStock + ent.AddedStock
	if ent.SoldQuantity+ent.RejectedQuantity > totalAvail {
		return fmt.Errorf("sold and rejected quantities cannot exceed total available stock of %d", totalAvail)
	}
	return nil
}

func (s *ProduceEntryService) CreateEntry(fID uuid.UUID, pID uuid.UUID, date time.Time, openStock, addStock, soldQty, rejQty, price int, notes string) (*models.ProduceEntry, error) {
	//verify product exists and belongs to this  farmer
	prod, err := s.pRepo.GetProductByID(pID)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}
	if prod.FarmerID != fID {
		return nil, fmt.Errorf("you do not have permission to use this this product")
	}

	ent := &models.ProduceEntry{
		ID:               uuid.New(),
		FarmerID:         fID,
		ProductID:        pID,
		EntryDate:        date,
		OpeningStock:     openStock,
		AddedStock:       addStock,
		SoldQuantity:     soldQty,
		RejectedQuantity: rejQty,
		PricePerUnit:     price,
		Notes:            notes,
		Version:          1,
		Deleted:          false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := validate(ent); err != nil {
		return nil, err
	}

	if err := s.eRepo.CreateProduceEntry(ent); err != nil {  //save to database
		return nil, fmt.Errorf("failed to create entry: %w", err)
	}

	compute(ent)
	return ent, nil
}

func (s *ProduceEntryService) GetEntries(fID uuid.UUID, start, end time.Time, pID *uuid.UUID) ([]*models.ProduceEntry, error) {
	entries, err := s.eRepo.GetEntriesByFarmerAndDateRange(fID, start, end, pID)
	if err != nil {
		return nil, fmt.Errorf("failed to get entries: %w", err)
	}
	for _, ent := range entries {
		compute(ent)
	}
	return entries, nil
}

//GetEntryByID returns a single entry by ID. Enforces ownership — a farmer can only view their own entries.
func (s *ProduceEntryService) GetEntryByID(id, farmerID uuid.UUID) (*models.ProduceEntry, error) {
	ent, err := s.eRepo.GetProduceEntryByID(id)
	if err != nil {
		return nil, fmt.Errorf("entry not found")
	}

	// Enforce ownership
	if ent.FarmerID != farmerID {
		return nil, fmt.Errorf("you do not have permission to view this entry")
	}

	compute(ent)

	return ent, nil
}

func (s *ProduceEntryService) UpdateEntry(id, fID uuid.UUID, pID uuid.UUID, date time.Time, openStock, addStock, soldQty, rejQty, price, version int, notes string) (*models.ProduceEntry, error) {
	ent, err := s.eRepo.GetProduceEntryByID(id)
	if err != nil {
		return nil, fmt.Errorf("entry not found")
	}
	if ent.FarmerID != fID {
		return nil, fmt.Errorf("you do not have permissionto update this entry")
	}
	ent.ProductID        = pID
	ent.EntryDate        = date
	ent.OpeningStock     = openStock
	ent.AddedStock       = addStock
	ent.SoldQuantity     = soldQty
	ent.RejectedQuantity = rejQty
	ent.PricePerUnit     = price
	ent.Notes            = notes
	ent.Version          = version
	ent.UpdatedAt        = time.Now()

	if err := validate(ent); err != nil {
		return nil, err
	}
	if err := s.eRepo.UpdateEntryWithVersionCheck(ent); err != nil {
		return nil, err
	}
	updated, err := s.eRepo.GetProduceEntryByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated entry: %w", err)
	}
	compute(updated)
	return updated, nil
}

func (s *ProduceEntryService) DeleteEntry(id,fID uuid.UUID) error {
	ent, err := s.eRepo.GetProduceEntryByID(id)
	if err != nil {
		return fmt.Errorf("entry not found")
	}

	if ent.FarmerID != fID {
		return fmt.Errorf("you do not have permission to delete this entry")
	}
	return s.eRepo.SoftDeleteEntry(id, fID)
}
