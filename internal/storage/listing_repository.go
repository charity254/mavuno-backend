package storage

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/mavuno/mavuno-backend/internal/models"
)

type ListingRepository struct {
	db *sql.DB
}

func NewListingRepository(db *sql.DB) *ListingRepository {
	return &ListingRepository{db: db}
}

func (r *ListingRepository) CreateListing(l *models.Listing) error {
	query := `
		INSERT INTO listings (id, farmer_id, title, description, price, quantity, unit_type, location, version, deleted, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.Exec(query,
		l.ID, l.FarmerID, l.Title, l.Description,
		l.Price, l.Quantity, l.UnitType, l.Location,
		l.Version, l.Deleted, l.CreatedAt, l.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create listing: %w", err)
	}
	return nil
}

func (r *ListingRepository) GetListingByID(id uuid.UUID) (*models.Listing, error) {
	query := `
		SELECT id, farmer_id, title, description, price, quantity, unit_type, location, version, deleted, created_at, updated_at
		FROM listings
		WHERE id = $1 AND deleted = false
	`
	l := &models.Listing{}
	err := r.db.QueryRow(query, id).Scan(
		&l.ID, &l.FarmerID, &l.Title, &l.Description,
		&l.Price, &l.Quantity, &l.UnitType, &l.Location,
		&l.Version, &l.Deleted, &l.CreatedAt, &l.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("listing not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}
	return l, nil
}

func (r *ListingRepository) GetAllActiveListings() ([]*models.Listing, error) {
	query := `
		SELECT id, farmer_id, title, description, price, quantity, unit_type, location, version, deleted, created_at, updated_at
		FROM listings
		WHERE deleted = false
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get listings: %w", err)
	}
	defer rows.Close()

	var listings []*models.Listing
	for rows.Next() {
		l := &models.Listing{}
		err := rows.Scan(
			&l.ID, &l.FarmerID, &l.Title, &l.Description,
			&l.Price, &l.Quantity, &l.UnitType, &l.Location,
			&l.Version, &l.Deleted, &l.CreatedAt, &l.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan listing: %w", err)
		}
		listings = append(listings, l)
	}
	return listings, nil
}

func (r *ListingRepository) GetListingsByFarmer(farmerID uuid.UUID) ([]*models.Listing, error) {
	query := `
		SELECT id, farmer_id, title, description, price, quantity, unit_type, location, version, deleted, created_at, updated_at
		FROM listings
		WHERE farmer_id = $1 AND deleted = false
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, farmerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get listings: %w", err)
	}
	defer rows.Close()

	var listings []*models.Listing
	for rows.Next() {
		l := &models.Listing{}
		err := rows.Scan(
			&l.ID, &l.FarmerID, &l.Title, &l.Description,
			&l.Price, &l.Quantity, &l.UnitType, &l.Location,
			&l.Version, &l.Deleted, &l.CreatedAt, &l.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan listing: %w", err)
		}
		listings = append(listings, l)
	}
	return listings, nil
}

func (r *ListingRepository) UpdateListingWithVersionCheck(l *models.Listing) error {
	query := `
		UPDATE listings
		SET title = $1, description = $2, price = $3, quantity = $4, unit_type = $5, location = $6, version = version + 1, updated_at = $7
		WHERE id = $8 AND version = $9 AND deleted = false
	`
	result, err := r.db.Exec(query,
		l.Title, l.Description, l.Price, l.Quantity,
		l.UnitType, l.Location, l.UpdatedAt, l.ID, l.Version,
	)
	if err != nil {
		return fmt.Errorf("failed to update listing: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check update result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("conflict: listing was updated by another session")
	}
	return nil
}

func (r *ListingRepository) SoftDeleteListing(id, farmerID uuid.UUID) error {
	query := `
		UPDATE listings
		SET deleted = true, updated_at = NOW()
		WHERE id = $1 AND farmer_id = $2 AND deleted = false
	`
	result, err := r.db.Exec(query, id, farmerID)
	if err != nil {
		return fmt.Errorf("failed to delete listing: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("listing not found or already deleted")
	}
	return nil
}