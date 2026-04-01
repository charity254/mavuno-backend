package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mavuno/mavuno-backend/internal/models"
	"github.com/mavuno/mavuno-backend/internal/storage"
)

type SupplyLocationService struct {
	repo *storage.SupplyLocationRepository
}

func NewSupplyLocationService(repo *storage.SupplyLocationRepository) *SupplyLocationService {
	return &SupplyLocationService{repo: repo}
}

func (s *SupplyLocationService) CreateSupplyLocation(farmerID uuid.UUID, name, contactPerson, phoneNumber, locationAddress, notes string) (*models.SupplyLocation, error) {
	// Validate required fields
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if contactPerson == "" {
		return nil, fmt.Errorf("contact person is required")
	}
	if phoneNumber == "" {
		return nil, fmt.Errorf("phone number is required")
	}
	if locationAddress == "" {
		return nil, fmt.Errorf("location address is required")
	}

	loc := &models.SupplyLocation{
		ID:              uuid.New(),
		FarmerID:        farmerID,
		Name:            name,
		ContactPerson:   contactPerson,
		PhoneNumber:     phoneNumber,
		LocationAddress: locationAddress,
		Notes:           notes,
		Version:         1,
		Deleted:         false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.repo.CreateSupplyLocation(loc); err != nil {
		return nil, fmt.Errorf("failed to create supply location: %w", err)
	}

	return loc, nil
}


func (s *SupplyLocationService) GetSupplyLocations(farmerID uuid.UUID) ([]*models.SupplyLocation, error) {
	locs, err := s.repo.GetSupplyLocationByFarmer(farmerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get supply locations: %w", err)
	}
	return locs, nil
}

func (s *SupplyLocationService) GetSupplyLocationByID(id, farmerID uuid.UUID) (*models.SupplyLocation, error) {
	loc, err := s.repo.GetSupplyLocationByID(id)
	if err != nil {
		return nil, fmt.Errorf("supply location not found")
	}

	// Enforce ownership
	if loc.FarmerID != farmerID {
		return nil, fmt.Errorf("you do not have permission to view this supply location")
	}

	return loc, nil
}

func (s *SupplyLocationService) UpdateSupplyLocation(id, farmerID uuid.UUID, name, contactPerson, phoneNumber, locationAddress, notes string, version int) (*models.SupplyLocation, error) {
	// Validate required fields
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if contactPerson == "" {
		return nil, fmt.Errorf("contact person is required")
	}
	if phoneNumber == "" {
		return nil, fmt.Errorf("phone number is required")
	}
	if locationAddress == "" {
		return nil, fmt.Errorf("location address is required")
	}

	// Fetch existing location to confirm ownership
	loc, err := s.repo.GetSupplyLocationByID(id)
	if err != nil {
		return nil, fmt.Errorf("supply location not found")
	}

	// Enforce ownership
	if loc.FarmerID != farmerID {
		return nil, fmt.Errorf("you do not have permission to update this supply location")
	}

	// Apply updates
	loc.Name            = name
	loc.ContactPerson   = contactPerson
	loc.PhoneNumber     = phoneNumber
	loc.LocationAddress = locationAddress
	loc.Notes           = notes
	loc.Version         = version
	loc.UpdatedAt       = time.Now()

	if err := s.repo.UpdateSupplyLocationWithVersionCheck(loc); err != nil {
		return nil, err
	}

	// Fetch updated location to get correct incremented version
	updated, err := s.repo.GetSupplyLocationByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated supply location: %w", err)
	}

	return updated, nil
}

func (s *SupplyLocationService) DeleteSupplyLocation(id, farmerID uuid.UUID) error {
	// Fetch location first to confirm ownership
	loc, err := s.repo.GetSupplyLocationByID(id)
	if err != nil {
		return fmt.Errorf("supply location not found")
	}

	// Enforce ownership
	if loc.FarmerID != farmerID {
		return fmt.Errorf("you do not have permission to delete this supply location")
	}

	return s.repo.SoftDeleteSupplyLocation(id, farmerID)
}