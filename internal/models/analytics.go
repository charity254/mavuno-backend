package models

// Each point is one day's total revenue.
type RevenuePoint struct {
	Date    string `json:"date"`    // format: YYYY-MM-DD
	Revenue int    `json:"revenue"` // total revenue in KES cents
}

// StockPoint represents a single data point on the stock chart.
// Each point is one day's remaining stock for a product.
type StockPoint struct {
	Date           string `json:"date"`
	RemainingStock int    `json:"remaining_stock"`
}

// RejectionPoint represents a single data point on the rejection chart.
// Each point is one day's total rejected quantity for a product.
type RejectionPoint struct {
	Date             string `json:"date"`
	RejectedQuantity int    `json:"rejected_quantity"`
}

// ProductSummary represents aggregated totals for a single product.
// Used in the product summary table on the analytics page.
type ProductSummary struct {
	ProductName    string `json:"product_name"`
	TotalSold      int    `json:"total_sold"`
	TotalRejected  int    `json:"total_rejected"`
	TotalRevenue   int    `json:"total_revenue"` // in KES cents
}

// PlannedVsActualPoint represents a single data point on the
// planned vs actual chart. Shows difference between what was
// agreed to deliver and what was actually sold.
type PlannedVsActualPoint struct {
	Date             string `json:"date"`
	PlannedQuantity  int    `json:"planned_quantity"`
	ActualSold       int    `json:"actual_sold"`
	Difference       int    `json:"difference"` // actual - planned
}