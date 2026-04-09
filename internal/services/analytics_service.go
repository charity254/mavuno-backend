package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mavuno/mavuno-backend/internal/models"
	"github.com/mavuno/mavuno-backend/internal/storage"
)

type AnalyticsService struct {
	repo *storage.AnalyticsRepository
}

func NewAnalyticsService(repo *storage.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

// validateDateRange checks that start and end dates are provided and that start is not after end.
func validateDateRange(start, end time.Time) error {
	if start.IsZero() || end.IsZero() {
		return fmt.Errorf("start and end dates are required")
	}
	if start.After(end) {
		return fmt.Errorf("start date cannot be after end date")
	}
	return nil
}

// GetRevenueTrend returns daily revenue totals for a farmer.
func (s *AnalyticsService) GetRevenueTrend(farmerID uuid.UUID, start, end time.Time) ([]models.RevenuePoint, error) {
	if err := validateDateRange(start, end); err != nil {
		return nil, err
	}
	points, err := s.repo.GetRevenueTrend(farmerID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue trend: %w", err)
	}
	// Return empty slice instead of nil so frontend gets [] not null
	if points == nil {
		points = []models.RevenuePoint{}
	}
	return points, nil
}

// GetStockTrend returns daily remaining stock for a specific product.
func (s *AnalyticsService) GetStockTrend(farmerID, productID uuid.UUID, start, end time.Time) ([]models.StockPoint, error) {
	if err := validateDateRange(start, end); err != nil {
		return nil, err
	}
	points, err := s.repo.GetStockTrend(farmerID, productID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock trend: %w", err)
	}
	if points == nil {
		points = []models.StockPoint{}
	}
	return points, nil
}

// GetRejectionTrend returns daily rejected quantities for a specific product.
func (s *AnalyticsService) GetRejectionTrend(farmerID, productID uuid.UUID, start, end time.Time) ([]models.RejectionPoint, error) {
	if err := validateDateRange(start, end); err != nil {
		return nil, err
	}
	points, err := s.repo.GetRejectionTrend(farmerID, productID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get rejection trend: %w", err)
	}
	if points == nil {
		points = []models.RejectionPoint{}
	}
	return points, nil
}

// GetProductSummary returns aggregated totals per product for a farmer.
func (s *AnalyticsService) GetProductSummary(farmerID uuid.UUID, start, end time.Time) ([]models.ProductSummary, error) {
	if err := validateDateRange(start, end); err != nil {
		return nil, err
	}
	summaries, err := s.repo.GetProductSummary(farmerID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get product summary: %w", err)
	}
	if summaries == nil {
		summaries = []models.ProductSummary{}
	}
	return summaries, nil
}

// GetPlannedVsActual compares agreed delivery quantities against actual sold.
func (s *AnalyticsService) GetPlannedVsActual(farmerID, productID uuid.UUID, start, end time.Time) ([]models.PlannedVsActualPoint, error) {
	if err := validateDateRange(start, end); err != nil {
		return nil, err
	}
	points, err := s.repo.GetPlannedVsActual(farmerID, productID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get planned vs actual: %w", err)
	}
	if points == nil {
		points = []models.PlannedVsActualPoint{}
	}
	return points, nil
}

// GetTodaySummary returns aggregated totals for today.
// Used by the dashboard summary cards.
func (s *AnalyticsService) GetTodaySummary(farmerID uuid.UUID) (*models.TodaySummary, error) {
	// Use today's date truncated to midnight so time component does not affect the query
	today := time.Now().Truncate(24 * time.Hour)
	summary, err := s.repo.GetTodaySummary(farmerID, today)
	if err != nil {
		return nil, fmt.Errorf("failed to get today summary: %w", err)
	}
	return summary, nil
}