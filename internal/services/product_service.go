package services

import (
	"fmt"
	"time"

	
    "github.com/google/uuid"
    "github.com/mavuno/mavuno-backend/internal/models"
    "github.com/mavuno/mavuno-backend/internal/storage"
)

type ProductService struct {
	productRepo *storage.ProductRepository
}

func NewProductRepository(productRepo *storage.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) CreateProduct(farmerID uuid.UUID, name, unitType, description string) (*models.Product, error) {
	if name == "" {
		return nil, fmt.Errorf("product name is required")
	}
	if unitType == "" {
		return nil, fmt.Errorf("unit type is required")
	}

	product := &models.Product{
		ID:          uuid.New(),   // generate a unique ID
        FarmerID:    farmerID,     // link the product to the farmer
        Name:        name,
        UnitType:    unitType,
        Description: description,
        Version:     1,            
        Deleted:     false,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
	}

	err := s.productRepo.CreateProduct(product)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	return product, nil
}

func(s *ProductService) GetProductsByFarmer(farmerID uuid.UUID) ([]*models.Product, error) {
	products, err := s.productRepo.GetProductByID(id)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}
	if product.FarmerID != farmerID {
		return nil, fmt.Errorf("you do not have permission to access this product")
	}
	return product, nil
}

func (s *ProductService) UpdateProduct(id, farmerID uuid.UUID, name, unitType, description string, version int) (*models.Product, error) {
	if name == "" {
		return nil, fmt.Errorf("product name is required")
	}
	if unitType == ""{
		return nil, fmt.Errorf("unit type is required")
	}
	product, err := s.productRepo.GetProductByID(id)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}
	if product.FarmerID != farmerID {
		return nil, fmt.Errorf("you do not have permission to update this product")
	}

	product.Name 		= name
	product.UnitType 	= unitType
	product.Description = description
	product.Version 	= version
	product.UpdatedAt 	= time.Now()

	err = s.productRepo.UpdateProductWithVersionCheck(product)
	if err != nil {
		if err.Error() == "conflict: product was updated by another session" {
			return nil, fmt.Errorf("conflict: product  was updated by another session")
		}
		return nil, fmt.Errorf("failed to update product: %w", err)
	}
	return product, nil
}

func (s *ProductService) DeleteProduct(id, farmerID uuid.UUID) error {
	err := s.productRepo.SoftDeleteProduct(id, farmerID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}