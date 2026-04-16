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