package models

type RevenuePoint struct {
	Date    string `json:"date"`    // format: YYYY-MM-DD
	Revenue int    `json:"revenue"` // total revenue in KES cents
}


type StockPoint struct {
	Date           string `json:"date"`
	RemainingStock int    `json:"remaining_stock"`
}


type RejectionPoint struct {
	Date             string `json:"date"`
	RejectedQuantity int    `json:"rejected_quantity"`
}


type ProductSummary struct {
	ProductName    string `json:"product_name"`
	TotalSold      int    `json:"total_sold"`
	TotalRejected  int    `json:"total_rejected"`
	TotalRevenue   int    `json:"total_revenue"` // in KES cents
}

type PlannedVsActualPoint struct {
	Date             string `json:"date"`
	PlannedQuantity  int    `json:"planned_quantity"`
	ActualSold       int    `json:"actual_sold"`
	Difference       int    `json:"difference"` // actual - planned
}

type TodaySummary struct {
    TotalAvailable int `json:"total_available"`
    TotalSold      int `json:"total_sold"`
    TotalRevenue   int `json:"total_revenue"`
    TotalRejected  int `json:"total_rejected"`
}