package storage

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
    "github.com/mavuno/mavuno-backend/internal/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) CreateProduct(p *models.Product) error {
	query := `
		INSERT INTO products (id, farmer_id, name, unit_type, description, version, deleted, created_at, updated)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(query,
		p.ID,
		p.FarmerID,
		p.Name,
		p.Description,
		p.Version,
		p.Deleted,
		p.CreatedAt,
		p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create product %w", err)
	}
	return nil
}

func (r *ProductRepository) GetProductByID(id uuid.UUID) (*models.Product, error) {
	query := `
		SELECT id, farmer_id, name, unit_type,  description, version, deleted, created_at, updated_at
		FROM products
		WHERE id = $1 AND deleted = false
	`
	p := &models.Product{}
	err := r.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.FarmerID,
		&p.Name,
		&p.UnitType,
		&p.Description,
		&p.Version,
		&p.Deleted,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return p, nil
}
func (r *ProductRepository) GetProductsByFarmer(farmerID uuid.UUID) ([]*models.Product, error) {
	query := `
		SELECT id, farmer_id, name, unit_type, description, version, deleted, created_at, updated_at
		FROM products
		WHERE farmer_id = $1 AND deleted = false
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, farmerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	defer rows.Close()


	var products []*models.Product
	for rows.Next() {
		p:= &models.Product{}
		err := rows.Scan(
			&p.ID,
			&p.FarmerID,
			&p.UnitType,
			&p.Description,
			&p.Version,
			&p.Deleted,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) UpdateProductWithVersionCheck(p *models.Product) error {
	query := `
		UPDATE products
		SET name = $1, unit_type = $2, description = $3, version = version + 1, updated_at = $4
		WHERE id = $5 AND version = $6 AND deleted = false
	`
	result, err := r.db.Exec(query,
		p.Name,
		p.UnitType,
		p.Description,
		p.UpdatedAt,
		p.ID,
		p.Version,
	)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check update result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("conflict: product was updated by another session")
	}
	return nil
}

func (r *ProductRepository) SoftDeleteProduct(id, farmerID uuid.UUID) error {
	query := `
		UPDATE products
		SET deleted = true, updated_at = NOW()
		WHERE id = $1 AND farmer_id = $2 AND deleted = false
	`
	result, err := r.db.Exec(query, id, farmerID)
	if err != nil {
		return fmt.Errorf("failed to delete product %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("product not found or already deleted")
	}
	return nil
}