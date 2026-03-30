package storage

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/mavuno/mavuno-backend/internal/models"
)

type SupplyLocationRepository struct {
	db *sql.DB
}

func NewSupplyLocationRepository(db *sql.DB) *SupplyLocationRepository {
	return &SupplyLocationRepository{db: db}
}

func (r *SupplyLocationRepository) CreateSupplyLocation(loc *models.SupplyLocation) error {
	query := `
		INSERT INTO supply_locations (id, farmer_id, name, contact_person, phone_number, location_address, notes, version, deleted, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Exec(query,
		loc.ID,
		loc.FarmerID,
		loc.Name,
		loc.ContactPerson,
		loc.PhoneNumber,
		loc.LocationAddress,
		loc.Notes,
		loc.Version,
		loc.Deleted,
		loc.CreatedAt,
		loc.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create supply location: %w", err)
	}
	return nil
}

func (r *SupplyLocationRepository) GetSupplyLocationByID(id uuid.UUID) (*models.SupplyLocation, error) {
		query := `
		SELECT id, farmer_id, name, contact_person, phone_number, location_address, notes, version, deleted, created_at, updated_at
		FROM supply_locations
		WHERE id = $1 AND deleted = false
		`
		loc := &models.SupplyLocation{}
		err := r.db.QueryRow(query, id).Scan(
			&loc.ID,
			&loc.FarmerID,
			&loc.Name,
			&loc.ContactPerson,
			&loc.PhoneNumber,
			&loc.LocationAddress,
			&loc.Notes,
			&loc.Version,
			&loc.Deleted,
			&loc.CreatedAt,
			&loc.UpdatedAt,
		)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("supply location not found")
		}
		if err != nil {
			return nil, fmt.Errorf("failed to get supply location: %w", err)
		}
		return loc, nil
}

func (r *SupplyLocationRepository) GetSupplyLocationByFarmer (farmerID uuid.UUID) ([]*models.SupplyLocation, error) {
	query := `
		SELECT id, farmer_id, name, contact_person, phone_number, location_address, notes, version, deleted, created_at, updated_at
		FROM supply_locations
		WHERE farmer_id = $1 AND deleted = false
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, farmerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get supply locations: %w", err)
	}
	defer rows.Close()

	var locations []*models.SupplyLocation
	for rows.Next() {
		loc := &models.SupplyLocation{}
		err := rows.Scan(
			&loc.ID,
			&loc.FarmerID,
			&loc.Name,
			&loc.ContactPerson,
			&loc.PhoneNumber,
			&loc.LocationAddress,
			&loc.Notes,
			&loc.Version,
			&loc.Deleted,
			&loc.CreatedAt,
			&loc.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan supply location: %w", err)
		}
		locations = append(locations, loc)
	}
	return locations, nil
}

func (r *SupplyLocationRepository) UpdateSupplyLocationWithVersionCheck(loc *models.SupplyLocation) error {
	query := `
		UPDATE supply_locations
		SET name = $1, contact_person = $2, phone_number = $3, location_address = $4, notes = $5, version = version + 1, updated_at = $6
		WHERE id = $7 AND version = $8 AND deleted = false
	`
	result, err := r.db.Exec(query,
		loc.Name,
		loc.ContactPerson,
		loc.PhoneNumber,
		loc.LocationAddress,
		loc.Notes,
		loc.UpdatedAt,
		loc.ID,
		loc.Version,
	)
	if err != nil {
		return fmt.Errorf("failed to update supply location: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check update result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("conflict: supply location was updated by another session")
	}
	return nil
}

func (r *SupplyLocationRepository) SoftDeleteSupplyLocation(id, farmerID uuid.UUID) error {
	query := `
		UPDATE supply_locations
		SET deleted = true, updated_at = NOW()
		WHERE id = $1 AND farmer_id = $2 AND deleted = false
	`
	result, err := r.db.Exec(query, id, farmerID)
	if err != nil {
		return fmt.Errorf("failed to delete supply location: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("supply location not found or already deleted")
	}
	return nil
}