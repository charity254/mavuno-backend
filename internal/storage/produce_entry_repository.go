package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mavuno/mavuno-backend/internal/models"
)

type ProduceEntryRepository struct {
	db *sql.DB
}

func NewProduceEntryRepository(db *sql.DB) *ProduceEntryRepository {
	return &ProduceEntryRepository{db: db}
}

func (r *ProduceEntryRepository) CreateProduceEntry(e *models.ProduceEntry) error {
	query := `
		INSERT INTO produce_entries (id, farmer_id, product-id, entry_date, opening-stock, added-stock, sold_quantity, rejected_quantity, price_per-unit, notes, version, deleted, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	_, err := r.db.Exec(query,
			e.ID,
		e.FarmerID,
		e.ProductID,
		e.EntryDate,
		e.OpeningStock,
		e.AddedStock,
		e.SoldQuantity,
		e.RejectedQuantity,
		e.PricePerUnit,
		e.Notes,
		e.Version,
		e.Deleted,
		e.CreatedAt,
		e.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create poduce entry: %w", err)
	}
	return nil
}
 
func (r *ProduceEntryRepository) GetProduceEntryByID(id uuid.UUID) (*models.ProduceEntry, error) {
	query := `
		SELECT id, farmer-id, product-id, entry-date, opening_stock, added_stock, sold_quantity, rejected-quantity, price_per_unit, notes, version, deleted, created_at, updated_at
		FROM produce_entries
		WHERE id = $1 AND deleted = false
	`
	e := &models.ProduceEntry{}
	err := r.db.QueryRow(query, id).Scan(
		&e.ID,
		&e.FarmerID,
		&e.ProductID,
		&e.EntryDate,
		&e.OpeningStock,
		&e.AddedStock,
		&e.SoldQuantity,
		&e.RejectedQuantity,
		&e.PricePerUnit,
		&e.Notes,
		&e.Version,
		&e.Deleted,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("entry not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get entry: %w", err)
	}
	return e, nil
}

func (r *ProduceEntryRepository) GetEntriesByFarmerAndDateRange(farmerID uuid.UUID, start, end time.Time, productID *uuid.UUID) ([]*models.ProduceEntry, error) { //fetches all entries for a farmer within a date range
	//qyuery is built dynamically based on whether a product filter is provided
	query := `
		SELECT id, farmer_id, product_id, entry_date, opening_stock, added_stock, sold_quantity, rejected_quantity, price_per_unit, notes, version, deleted, created_at, updated_at
		FROM produce_entries
		WHERE farmer_id = $1 AND deleted = false
	`
	args := []interface{}{farmerID}

	//add date range filter if provided
	if !start.IsZero() && !end.IsZero() {
		query += fmt.Sprintf(" AND entry_date BETWEEN $%d AND $%d", len(args)+1, len(args)+2)
		args = append(args, start, end)
	}
	//add product filter if provided
	if productID != nil {
		query += fmt.Sprintf(" AND prdoduct_id = $%d", len(args)+1)
		args = append(args, *productID)
	}
	query += " ORDER BY entry_date DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get entries: %w", err)
	}
	defer rows.Close()

	var entries []*models.ProduceEntry
	for rows.Next() {
		e := &models.ProduceEntry{}
		err := rows.Scan(
			&e.ID,
			&e.FarmerID,
			&e.ProductID,
			&e.EntryDate,
			&e.OpeningStock,
			&e.AddedStock,
			&e.SoldQuantity,
			&e.RejectedQuantity,
			&e.PricePerUnit,
			&e.Notes,
			&e.Version,
			&e.Deleted,
			&e.CreatedAt,
			&e.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func (r *ProduceEntryRepository) UpdateEntryWithVersionCheck(e *models.ProduceEntry) error { //updates a produce entry only if the version matches.
	query := `
		UPDATE produce_entries
		SET product_id = $1, entry_date = $2, opening_stock = $3, added_stock = $4, sold_quantity = $5, rejected_quantity = $6, price_per_unit = $7, notes = $8, version = version + 1, updated_at = $9
		WHERE id = $10 AND version = $11 AND deleted = false
	`
	result, err := r.db.Exec(query, 
		e.ProductID,
		e.EntryDate,
		e.OpeningStock,
		e.AddedStock,
		e.SoldQuantity,
		e.RejectedQuantity,
		e.PricePerUnit,
		e.Notes,
		e.UpdatedAt,
		e.ID,
		e.Version,
	)
	if err != nil {
		return fmt.Errorf("failed to update entry: %w", err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return fmt.Errorf("failed to check update result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("conflict: entry was updated by another session")
	}
	return nil
}

func (r *ProduceEntryRepository) SoftDeleteEntry(id, farmerID uuid.UUID) error {
	query := `
		UPDATE produce_entries
		SET deleted = true, updated_at =NOW()
		WHERE id = $1 AND farmer_id = $2 AND deleted = false
	`

	result, err := r.db.Exec(query, id, farmerID)
	if err != nil {
		return fmt.Errorf("failed to delete entry: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("entry not found or already deleted")
	}
	return nil
}
