package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mavuno/mavuno-backend/internal/services"
	"github.com/mavuno/mavuno-backend/internal/utils"
)

type ReportHandler struct {
	svc *services.ReportService
}

func NewReportHandler(svc *services.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

// parsReportParams reads and validates the common query params for reports
func parseReportParams(r *http.Request) (start, end time.Time, productID *uuid.UUID, err error) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	if startStr == "" || endStr == "" {
		err = fmt.Errorf("start and end dates are required")
		return
	}

	start, err = time.Parse("2006-01-02", startStr)
	if err != nil {
		err = fmt.Errorf("invalid start date — use YYYY-MM-DD")
		return
	}

	end, err = time.Parse("2006-01-02", endStr)
	if err != nil {
		err = fmt.Errorf("invalid end date — use YYYY-MM-DD")
		return
	}

	productIDStr := r.URL.Query().Get("product_id")
	if productIDStr != "" {
		pid, parseErr := uuid.Parse(productIDStr)
		if parseErr != nil {
			err = fmt.Errorf("invalid product ID")
			return
		}
		productID = &pid
	}

	return
}

// GetReport handles GET /api/reports/export
// Returns JSON data for displaying the report in the web app
func (h *ReportHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	start, end, productID, err := parseReportParams(r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	report, err := h.svc.GenerateReport(farmerID, start, end, productID)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, report)
}

// DownloadCSV handles GET /api/reports/download
// Returns a CSV file for downloading
func (h *ReportHandler) DownloadCSV(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	start, end, productID, err := parseReportParams(r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	csvData, err := h.svc.GenerateCSVReport(farmerID, start, end, productID)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	// Set headers to trigger file download in the browser
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=mavuno_report.csv")
	w.WriteHeader(http.StatusOK)
	w.Write(csvData)
}