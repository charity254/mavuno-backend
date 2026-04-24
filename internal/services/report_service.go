package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mavuno/mavuno-backend/internal/storage"
)

type ReportService struct {
	entryRepo   *storage.ProduceEntryRepository
	productRepo *storage.ProductRepository
}

// ReportRow represents a single row in the report.
type ReportRow struct {
	Date         string  `json:"date"`
	ProductName  string  `json:"product_name"`
	OpeningStock int     `json:"opening_stock"`
	AddedStock   int     `json:"added_stock"`
	TotalAvail   int     `json:"total_available"`
	Sold         int     `json:"sold"`
	Rejected     int     `json:"rejected"`
	Remaining    int     `json:"remaining"`
	PricePerUnit float64 `json:"price_per_unit_kes"`
	Revenue      float64 `json:"revenue_kes"`
}

// ReportSummary wraps the rows and totals for the JSON response.
type ReportSummary struct {
	Rows         []ReportRow `json:"rows"`
	TotalSold    int         `json:"total_sold"`
	TotalRejected int        `json:"total_rejected"`
	TotalRevenue float64     `json:"total_revenue_kes"`
}

func NewReportService(entryRepo *storage.ProduceEntryRepository, productRepo *storage.ProductRepository) *ReportService {
    return &ReportService{
        entryRepo:   entryRepo,
        productRepo: productRepo,
    }
}

func (s *ReportService) GenerateCSVReport(farmerID uuid.UUID, start, end time.Time, productID *uuid.UUID) ([]byte, error) {
	if start.IsZero() || end.IsZero() {
		return nil, fmt.Errorf("start and end dates are required")
	}
	if start.After(end) {
		return nil, fmt.Errorf("start date cannot be after end date")
	}

	// Fetch entries for the date range
	entries, err := s.entryRepo.GetEntriesByFarmerAndDateRange(farmerID, start, end, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch entries: %w", err)
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write CSV headers
	headers := []string{
		"Date",
		"Product",
		"Opening Stock",
		"Added Stock",
		"Total Available",
		"Sold",
		"Rejected",
		"Remaining",
		"Price Per Unit (KES)",
		"Revenue (KES)",
	}
	if err := writer.Write(headers); err != nil {
		return nil, fmt.Errorf("failed to write CSV headers: %w", err)
	}

	// Track totals for the summary row at the bottom
	var totalSold, totalRejected, totalRevenue int

	// Write one row per entry
	for _, entry := range entries {
		// Fetch product name for this entry
		product, err := s.productRepo.GetProductByID(entry.ProductID)
		if err != nil {
			product = nil
		}

		productName := "Unknown"
		if product != nil {
			productName = product.Name
		}

		// Compute derived fields
		totalAvail := entry.OpeningStock + entry.AddedStock
		remaining := totalAvail - entry.SoldQuantity - entry.RejectedQuantity
		revenue := entry.SoldQuantity * entry.PricePerUnit

		// Convert KES cents to KES for display
		priceKES := fmt.Sprintf("%.2f", float64(entry.PricePerUnit)/100)
		revenueKES := fmt.Sprintf("%.2f", float64(revenue)/100)

		row := []string{
			entry.EntryDate.Format("2006-01-02"),
			productName,
			fmt.Sprintf("%d", entry.OpeningStock),
			fmt.Sprintf("%d", entry.AddedStock),
			fmt.Sprintf("%d", totalAvail),
			fmt.Sprintf("%d", entry.SoldQuantity),
			fmt.Sprintf("%d", entry.RejectedQuantity),
			fmt.Sprintf("%d", remaining),
			priceKES,
			revenueKES,
		}

		if err := writer.Write(row); err != nil {
			return nil, fmt.Errorf("failed to write CSV row: %w", err)
		}

		totalSold += entry.SoldQuantity
		totalRejected += entry.RejectedQuantity
		totalRevenue += revenue
	}

	// Write totals row at the bottom
	totalsRow := []string{
		"TOTALS",
		"",
		"",
		"",
		"",
		fmt.Sprintf("%d", totalSold),
		fmt.Sprintf("%d", totalRejected),
		"",
		"",
		fmt.Sprintf("%.2f", float64(totalRevenue)/100),
	}
	if err := writer.Write(totalsRow); err != nil {
		return nil, fmt.Errorf("failed to write totals row: %w", err)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("failed to flush CSV: %w", err)
	}

	return buf.Bytes(), nil
}

// GenerateReport returns structured JSON data for displaying in the web app.
func (s *ReportService) GenerateReport(farmerID uuid.UUID, start, end time.Time, productID *uuid.UUID) (*ReportSummary, error) {
	if start.IsZero() || end.IsZero() {
		return nil, fmt.Errorf("start and end dates are required")
	}
	if start.After(end) {
		return nil, fmt.Errorf("start date cannot be after end date")
	}

	entries, err := s.entryRepo.GetEntriesByFarmerAndDateRange(farmerID, start, end, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch entries: %w", err)
	}

	var rows []ReportRow
	var totalSold, totalRejected int
	var totalRevenue float64

	for _, entry := range entries {
		product, err := s.productRepo.GetProductByID(entry.ProductID)
		productName := "Unknown"
		if err == nil && product != nil {
			productName = product.Name
		}

		totalAvail := entry.OpeningStock + entry.AddedStock
		remaining := totalAvail - entry.SoldQuantity - entry.RejectedQuantity
		revenue := float64(entry.SoldQuantity*entry.PricePerUnit) / 100

		rows = append(rows, ReportRow{
			Date:         entry.EntryDate.Format("2006-01-02"),
			ProductName:  productName,
			OpeningStock: entry.OpeningStock,
			AddedStock:   entry.AddedStock,
			TotalAvail:   totalAvail,
			Sold:         entry.SoldQuantity,
			Rejected:     entry.RejectedQuantity,
			Remaining:    remaining,
			PricePerUnit: float64(entry.PricePerUnit) / 100,
			Revenue:      revenue,
		})

		totalSold += entry.SoldQuantity
		totalRejected += entry.RejectedQuantity
		totalRevenue += revenue
	}

	if rows == nil {
		rows = []ReportRow{}
	}

	return &ReportSummary{
		Rows:          rows,
		TotalSold:     totalSold,
		TotalRejected: totalRejected,
		TotalRevenue:  totalRevenue,
	}, nil
}