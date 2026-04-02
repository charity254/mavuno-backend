package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mavuno/mavuno-backend/internal/models"
	"github.com/mavuno/mavuno-backend/internal/storage"
)

type SupplyAgreementService struct {
	repo     *storage.SupplyAgreementRepository
	locRepo  *storage.SupplyLocationRepository
	prodRepo *storage.ProductRepository
}

func NewSupplyAgreementService(repo *storage.SupplyAgreementRepository, locRepo *storage.SupplyLocationRepository, prodRepo *storage.ProductRepository) *SupplyAgreementService {
	return &SupplyAgreementService{
		repo:     repo,
		locRepo:  locRepo,
		prodRepo: prodRepo,
	}
}

var validDays = map[string]bool{
	"Monday":    true,
	"Tuesday":   true,
	"Wednesday": true,
	"Thursday":  true,
	"Friday":    true,
	"Saturday":  true,
	"Sunday":    true,
}

func validateDeliveryDays(days []string) error {
	if len(days) == 0 {
		return fmt.Errorf("at least one delivery day is required")
	}
	for _, day := range days {
		if !validDays[day] {
			return fmt.Errorf("invalid delivery day: %s — must be a full weekday name e.g. Monday", day)
		}
	}
	return nil
}

func (s *SupplyAgreementService) CreateSupplyAgreement(farmerID, productID, locationID uuid.UUID, qtyPerDelivery, pricePerUnit int, deliveryDays []string, deliveryNotes string) (*models.SupplyAgreement, error) {
	// Validate quantity
	if qtyPerDelivery <= 0 {
		return nil, fmt.Errorf("quantity per delivery must be greater than 0")
	}

	// Validate delivery days
	if err := validateDeliveryDays(deliveryDays); err != nil {
		return nil, err
	}

	// Verify product belongs to farmer
	prod, err := s.prodRepo.GetProductByID(productID)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}
	if prod.FarmerID != farmerID {
		return nil, fmt.Errorf("you do not have permission to use this product")
	}

	// Verify location belongs to farmer
	loc, err := s.locRepo.GetSupplyLocationByID(locationID)
	if err != nil {
		return nil, fmt.Errorf("supply location not found")
	}
	if loc.FarmerID != farmerID {
		return nil, fmt.Errorf("you do not have permission to use this supply location")
	}

	sa := &models.SupplyAgreement{
		ID:               uuid.New(),
		FarmerID:         farmerID,
		ProductID:        productID,
		SupplyLocationID: locationID,
		QtyPerDelivery:   qtyPerDelivery,
		PricePerUnit:     pricePerUnit,
		DeliveryDays:     deliveryDays,
		DeliveryNotes:    deliveryNotes,
		Active:           true,
		Version:          1,
		Deleted:          false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.repo.CreateSupplyAgreement(sa); err != nil {
		return nil, fmt.Errorf("failed to create supply agreement: %w", err)
	}

	return sa, nil
}

func (s *SupplyAgreementService) GetSupplyAgreements(farmerID uuid.UUID) ([]*models.SupplyAgreement, error) {
	agreements, err := s.repo.GetAgreementsByFarmer(farmerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get supply agreements: %w", err)
	}
	return agreements, nil
}


//to generate dashboard reminders.
func (s *SupplyAgreementService) GetActiveSupplyAgreements(farmerID uuid.UUID) ([]*models.SupplyAgreement, error) {
	agreements, err := s.repo.GetActiveAgreementsByFarmer(farmerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active supply agreements: %w", err)
	}
	return agreements, nil
}

func (s *SupplyAgreementService) GetSupplyAgreementByID(id, farmerID uuid.UUID) (*models.SupplyAgreement, error) {
	sa, err := s.repo.GetSupplyAgreementByID(id)
	if err != nil {
		return nil, fmt.Errorf("supply agreement not found")
	}

	// Enforce ownership
	if sa.FarmerID != farmerID {
		return nil, fmt.Errorf("you do not have permission to view this supply agreement")
	}

	return sa, nil
}

func (s *SupplyAgreementService) UpdateSupplyAgreement(id, farmerID, productID, locationID uuid.UUID, qtyPerDelivery, pricePerUnit int, deliveryDays []string, deliveryNotes string, active bool, version int) (*models.SupplyAgreement, error) {
	// Validate quantity
	if qtyPerDelivery <= 0 {
		return nil, fmt.Errorf("quantity per delivery must be greater than 0")
	}

	// Validate delivery days
	if err := validateDeliveryDays(deliveryDays); err != nil {
		return nil, err
	}

	// Fetch existing agreement to confirm ownership
	sa, err := s.repo.GetSupplyAgreementByID(id)
	if err != nil {
		return nil, fmt.Errorf("supply agreement not found")
	}

	// Enforce ownership
	if sa.FarmerID != farmerID {
		return nil, fmt.Errorf("you do not have permission to update this supply agreement")
	}

	// Apply updates
	sa.ProductID        = productID
	sa.SupplyLocationID = locationID
	sa.QtyPerDelivery   = qtyPerDelivery
	sa.PricePerUnit     = pricePerUnit
	sa.DeliveryDays     = deliveryDays
	sa.DeliveryNotes    = deliveryNotes
	sa.Active           = active
	sa.Version          = version
	sa.UpdatedAt        = time.Now()

	if err := s.repo.UpdateAgreementWithVersionCheck(sa); err != nil {
		return nil, err
	}

	// Fetch updated agreement for incremented version
	updated, err := s.repo.GetSupplyAgreementByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated supply agreement: %w", err)
	}

	return updated, nil
}

func (s *SupplyAgreementService) DeleteSupplyAgreement(id, farmerID uuid.UUID) error {
	sa, err := s.repo.GetSupplyAgreementByID(id)
	if err != nil {
		return fmt.Errorf("supply agreement not found")
	}

	// Enforce ownership
	if sa.FarmerID != farmerID {
		return fmt.Errorf("you do not have permission to delete this supply agreement")
	}

	return s.repo.SoftDeleteAgreement(id, farmerID)
}