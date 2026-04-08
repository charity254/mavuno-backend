package storage 

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mavuno/mavuno-backend/internal/models"
)

type AnalyticsRepository struct { //read only — no inserts or updates.
	db *sql.DB
}

func (r *AnalyticsRepository) GetRevenueTrend(farmerID uuid.UUID, start, end time.Time) ([]models.RevenuePoint, error) {
	query := `
		SELECT 
			entry_date::TEXT,
			SUM(sold_quantity * price_per_unit) AS revenue
		FROM produce_entries
		WHERE farmer_id = $1
			AND entry_date BETWEEN $2 AND $3
			AND deleted = false
		GROUP BY entry_date
		ORDER BY entry_date ASC
	`
	rows, err := r.db.Query(query, farmerID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue trend: %w", err)
	}
	defer rows.Close()

	var points []models.RevenuePoint
	for rows.Next() {
		var pt models.RevenuePoint
		if err := rows.Scan(&pt.Date, &pt.Revenue); err != nil {
			return nil, fmt.Errorf("failed to scan revenue point: %w", err)
		}
		points = append(points, pt)
	}
	return points, nil
}

func (r *AnalyticsRepository) GetStockTrend(farmerID, productID uuid.UUID, start, end time.Time) ([]models.StockPoint, error) {
	query := `
		SELECT
			entry_date::TEXT,
			SUM((opening_stock + added_stock) - sold_quantity - rejected_quantity) AS remaining_stock
		FROM produce_entries
		WHERE farmer_id = $1
			AND product_id = $2
			AND entry_date BETWEEN $3 AND $4
			AND deleted = false
		GROUP BY entry_date
		ORDER BY entry_date ASC
	`
	rows, err := r.db.Query(query, farmerID, productID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock trend: %w", err)
	}
	defer rows.Close()

	var points []models.StockPoint
	for rows.Next() {
		var pt models.StockPoint
		if err := rows.Scan(&pt.Date, &pt.RemainingStock); err != nil {
			return nil, fmt.Errorf("failed to scan stock point: %w", err)
		}
		points = append(points, pt)
	}
	return points, nil
}

func (r *AnalyticsRepository) GetRejectionTrend(farmerID, productID uuid.UUID, start, end time.Time) ([]models.RejectionPoint, error) {  //returns daily rejected quantities for a specific product.
	query := `
		SELECT
			entry_date::TEXT,
			SUM(rejected_quantity) AS rejected_quantity
		FROM produce_entries
		WHERE farmer_id = $1
			AND product_id = $2
			AND entry_date BETWEEN $3 AND $4
			AND deleted = false
		GROUP BY entry_date
		ORDER BY entry_date ASC
	`
	rows, err := r.db.Query(query, farmerID, productID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get rejection trend: %w", err)
	}
	defer rows.Close()

	var points []models.RejectionPoint
	for rows.Next() {
		var pt models.RejectionPoint
		if err := rows.Scan(&pt.Date, &pt.RejectedQuantity); err != nil {
			return nil, fmt.Errorf("failed to scan rejection point: %w", err)
		}
		points = append(points, pt)
	}
	return points, nil
}

func (r *AnalyticsRepository) GetProductSummary(farmerID uuid.UUID, start, end time.Time) ([]models.ProductSummary, error) { //returns aggregated totals per product for a farmer.
	query := `
		SELECT
			p.name AS product_name,
			SUM(e.sold_quantity) AS total_sold,
			SUM(e.rejected_quantity) AS total_rejected,
			SUM(e.sold_quantity * e.price_per_unit) AS total_revenue
		FROM produce_entries e
		JOIN products p ON p.id = e.product_id
		WHERE e.farmer_id = $1
			AND e.entry_date BETWEEN $2 AND $3
			AND e.deleted = false
			AND p.deleted = false
		GROUP BY p.name
		ORDER BY total_revenue DESC
	`
	rows, err := r.db.Query(query, farmerID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get product summary: %w", err)
	}
	defer rows.Close()

	var summaries []models.ProductSummary
	for rows.Next() {
		var s models.ProductSummary
		if err := rows.Scan(&s.ProductName, &s.TotalSold, &s.TotalRejected, &s.TotalRevenue); err != nil {
			return nil, fmt.Errorf("failed to scan product summary: %w", err)
		}
		summaries = append(summaries, s)
	}
	return summaries, nil
}

func (r *AnalyticsRepository) GetPlannedVsActual(farmerID, productID uuid.UUID, start, end time.Time) ([]models.PlannedVsActualPoint, error) {  // compares agreed delivery quantities against actual sold quantities.
	query := `
		SELECT
			e.entry_date::TEXT,
			COALESCE(SUM(sa.quantity_per_delivery), 0) AS planned_quantity,
			COALESCE(SUM(e.sold_quantity), 0) AS actual_sold,
			COALESCE(SUM(e.sold_quantity), 0) - COALESCE(SUM(sa.quantity_per_delivery), 0) AS difference
		FROM produce_entries e
		LEFT JOIN supply_agreements sa 
			ON sa.product_id = e.product_id 
			AND sa.farmer_id = e.farmer_id
			AND sa.active = true
			AND sa.deleted = false
		WHERE e.farmer_id = $1
			AND e.product_id = $2
			AND e.entry_date BETWEEN $3 AND $4
			AND e.deleted = false
		GROUP BY e.entry_date
		ORDER BY e.entry_date ASC
	`
	rows, err := r.db.Query(query, farmerID, productID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get planned vs actual: %w", err)
	}
	defer rows.Close()

	var points []models.PlannedVsActualPoint
	for rows.Next() {
		var pt models.PlannedVsActualPoint
		if err := rows.Scan(&pt.Date, &pt.PlannedQuantity, &pt.ActualSold, &pt.Difference); err != nil {
			return nil, fmt.Errorf("failed to scan planned vs actual point: %w", err)
		}
		points = append(points, pt)
	}
	return points, nil
}

func (r *AnalyticsRepository) GetTodaySummary(farmerID uuid.UUID, today time.Time) (*models.TodaySummary, error) {
	query := `
		SELECT
			COALESCE(SUM(opening_stock + added_stock), 0) AS total_available,
			COALESCE(SUM(sold_quantity), 0) AS total_sold,
			COALESCE(SUM(sold_quantity * price_per_unit), 0) AS total_revenue,
			COALESCE(SUM(rejected_quantity), 0) AS total_rejected
		FROM produce_entries
		WHERE farmer_id = $1
			AND entry_date = $2
			AND deleted = false
	`
	summary := &models.TodaySummary{}
	err := r.db.QueryRow(query, farmerID, today).Scan(
		&summary.TotalAvailable,
		&summary.TotalSold,
		&summary.TotalRevenue,
		&summary.TotalRejected,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get today summary: %w", err)
	}
	return summary, nil
}