package storage

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/mavuno/mavuno-backend/internal/models"
)


type SupplyAgreementRepository struct {
	db *sql.DB
}

func NewSupplyAgreementRepository(db *sql.DB) *SupplyAgreementRepository {
	return &SupplyAgreementRepository{db: db}
}

func (r *SupplyAgreementRepository) CreateSupplyAgreement(sa *models.SupplyAgreement) error {
	query := `
		INSERT INTO supply_agreements (id, farmer_id, product_id, supply_location_id, quantity_per_delivery, price_per_unit, delivery_days, delivery_notes, active, version, deleted, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.Exec(query,
		sa.ID,
		sa.FarmerID,
		sa.ProductID,
		sa.SupplyLocationID,
		sa.QtyPerDelivery,
		sa.PricePerUnit,
		pq.Array(sa.DeliveryDays), // pq.Array converts Go slice to PostgreSQL array
		sa.DeliveryNotes,
		sa.Active,
		sa.Version,
		sa.Deleted,
		sa.CreatedAt,
		sa.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create supply agreement: %w", err)
	}
	return nil
}

func (r *SupplyAgreementRepository) GetSupplyAgreementByID(id uuid.UUID) (*models.SupplyAgreement, error) {
	query := `
		SELECT id, farmer_id, product_id, supply_location_id, quantity_per_delivery, price_per_unit, delivery_days, delivery_notes, active, version, deleted, created_at, updated_at
		FROM supply_agreements
		WHERE id = $1 AND deleted = false
	`
	sa := &models.SupplyAgreement{}
	err := r.db.QueryRow(query, id).Scan(
		&sa.ID,
		&sa.FarmerID,
		&sa.ProductID,
		&sa.SupplyLocationID,
		&sa.QtyPerDelivery,
		&sa.PricePerUnit,
		pq.Array(&sa.DeliveryDays), // pq.Array converts PostgreSQL array back to Go slice
		&sa.DeliveryNotes,
		&sa.Active,
		&sa.Version,
		&sa.Deleted,
		&sa.CreatedAt,
		&sa.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("supply agreement not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get supply agreement: %w", err)
	}
	return sa, nil
}

func (r *SupplyAgreementRepository) GetAgreementsByFarmer(farmerID uuid.UUID) ([]*models.SupplyAgreement, error) {
	query := `
		SELECT id, farmer_id, product_id, supply_location_id, quantity_per_delivery, price_per_unit, delivery_days, delivery_notes, active, version, deleted, created_at, updated_at
		FROM supply_agreements
		WHERE farmer_id = $1 AND deleted = false
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, farmerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get supply agreements: %w", err)
	}
	defer rows.Close()

	var agreements []*models.SupplyAgreement
	for rows.Next() {
		sa := &models.SupplyAgreement{}
		err := rows.Scan(
			&sa.ID,
			&sa.FarmerID,
			&sa.ProductID,
			&sa.SupplyLocationID,
			&sa.QtyPerDelivery,
			&sa.PricePerUnit,
			pq.Array(&sa.DeliveryDays),
			&sa.DeliveryNotes,
			&sa.Active,
			&sa.Version,
			&sa.Deleted,
			&sa.CreatedAt,
			&sa.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan supply agreement: %w", err)
		}
		agreements = append(agreements, sa)
	}
	return agreements, nil
}


//  to generate dashboard reminders.
func (r *SupplyAgreementRepository) GetActiveAgreementsByFarmer(farmerID uuid.UUID) ([]*models.SupplyAgreement, error) {
	query := `
		SELECT id, farmer_id, product_id, supply_location_id, quantity_per_delivery, price_per_unit, delivery_days, delivery_notes, active, version, deleted, created_at, updated_at
		FROM supply_agreements
		WHERE farmer_id = $1 AND active = true AND deleted = false
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, farmerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active supply agreements: %w", err)
	}
	defer rows.Close()

	var agreements []*models.SupplyAgreement
	for rows.Next() {
		sa := &models.SupplyAgreement{}
		err := rows.Scan(
			&sa.ID,
			&sa.FarmerID,
			&sa.ProductID,
			&sa.SupplyLocationID,
			&sa.QtyPerDelivery,
			&sa.PricePerUnit,
			pq.Array(&sa.DeliveryDays),
			&sa.DeliveryNotes,
			&sa.Active,
			&sa.Version,
			&sa.Deleted,
			&sa.CreatedAt,
			&sa.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan active supply agreement: %w", err)
		}
		agreements = append(agreements, sa)
	}
	return agreements, nil
}

func (r *SupplyAgreementRepository) UpdateAgreementWithVersionCheck(sa *models.SupplyAgreement) error {
	query := `
		UPDATE supply_agreements
		SET product_id = $1, supply_location_id = $2, quantity_per_delivery = $3, price_per_unit = $4, delivery_days = $5, delivery_notes = $6, active = $7, version = version + 1, updated_at = $8
		WHERE id = $9 AND version = $10 AND deleted = false
	`
	result, err := r.db.Exec(query,
		sa.ProductID,
		sa.SupplyLocationID,
		sa.QtyPerDelivery,
		sa.PricePerUnit,
		pq.Array(sa.DeliveryDays),
		sa.DeliveryNotes,
		sa.Active,
		sa.UpdatedAt,
		sa.ID,
		sa.Version,
	)
	if err != nil {
		return fmt.Errorf("failed to update supply agreement: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check update result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("conflict: supply agreement was updated by another session")
	}
	return nil
}

func (r *SupplyAgreementRepository) SoftDeleteAgreement(id, farmerID uuid.UUID) error {
	query := `
		UPDATE supply_agreements
		SET deleted = true, updated_at = NOW()
		WHERE id = $1 AND farmer_id = $2 AND deleted = false
	`
	result, err := r.db.Exec(query, id, farmerID)
	if err != nil {
		return fmt.Errorf("failed to delete supply agreement: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("supply agreement not found or already deleted")
	}
	return nil
}