package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mavuno/mavuno-backend/internal/models"
	"github.com/mavuno/mavuno-backend/internal/storage"
)

type ListingService struct {
	repo *storage.ListingRepository
}

func NewListingService(repo *storage.ListingRepository) *ListingService {
	return &ListingService{repo: repo}
}

func (s *ListingService) CreateListing(farmerID uuid.UUID, title, description string, price, quantity int, unitType, location string) (*models.Listing, error) {
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if unitType == "" {
		return nil, fmt.Errorf("unit type is required")
	}
	if location == "" {
		return nil, fmt.Errorf("location is required")
	}
	if price < 0 {
		return nil, fmt.Errorf("price cannot be negative")
	}
	if quantity < 0 {
		return nil, fmt.Errorf("quantity cannot be negative")
	}

	l := &models.Listing{
		ID:          uuid.New(),
		FarmerID:    farmerID,
		Title:       title,
		Description: description,
		Price:       price,
		Quantity:    quantity,
		UnitType:    unitType,
		Location:    location,
		Version:     1,
		Deleted:     false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreateListing(l); err != nil {
		return nil, fmt.Errorf("failed to create listing: %w", err)
	}

	return l, nil
}

func (s *ListingService) GetAllListings() ([]*models.Listing, error) {
	listings, err := s.repo.GetAllActiveListings()
	if err != nil {
		return nil, fmt.Errorf("failed to get listings: %w", err)
	}
	if listings == nil {
		listings = []*models.Listing{}
	}
	return listings, nil
}

func (s *ListingService) GetListingByID(id uuid.UUID) (*models.Listing, error) {
	listing, err := s.repo.GetListingByID(id)
	if err != nil {
		return nil, fmt.Errorf("listing not found")
	}
	return listing, nil
}

func (s *ListingService) UpdateListing(id, farmerID uuid.UUID, title, description string, price, quantity int, unitType, location string, version int) (*models.Listing, error) {
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if unitType == "" {
		return nil, fmt.Errorf("unit type is required")
	}
	if location == "" {
		return nil, fmt.Errorf("location is required")
	}

	listing, err := s.repo.GetListingByID(id)
	if err != nil {
		return nil, fmt.Errorf("listing not found")
	}

	if listing.FarmerID != farmerID {
		return nil, fmt.Errorf("you do not have permission to update this listing")
	}

	listing.Title       = title
	listing.Description = description
	listing.Price       = price
	listing.Quantity    = quantity
	listing.UnitType    = unitType
	listing.Location    = location
	listing.Version     = version
	listing.UpdatedAt   = time.Now()

	if err := s.repo.UpdateListingWithVersionCheck(listing); err != nil {
		return nil, err
	}

	updated, err := s.repo.GetListingByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated listing: %w", err)
	}

	return updated, nil
}

func (s *ListingService) DeleteListing(id, farmerID uuid.UUID) error {
	listing, err := s.repo.GetListingByID(id)
	if err != nil {
		return fmt.Errorf("listing not found")
	}

	if listing.FarmerID != farmerID {
		return fmt.Errorf("you do not have permission to delete this listing")
	}

	return s.repo.SoftDeleteListing(id, farmerID)
}