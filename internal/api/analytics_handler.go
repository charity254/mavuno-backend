package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mavuno/mavuno-backend/internal/services"
	"github.com/mavuno/mavuno-backend/internal/utils"
)

type AnalyticsHandler struct {
	svc *services.AnalyticsService
}

func NewAnalyticsHandler(svc *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{svc: svc}
}

func parseAnalyticsParams(r *http.Request) (start, end time.Time, productID *uuid.UUID, err error) {
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	if startStr != "" {
		start, err = time.Parse("2006-01-02", startStr)
		if err != nil {
			err = fmt.Errorf("invalid start date — use YYYY-MM-DD")
			return
		}
	}

	if endStr != "" {
		end, err = time.Parse("2006-01-02", endStr)
		if err != nil {
			err = fmt.Errorf("invalid end date — use YYYY-MM-DD")
			return
		}
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

func (h *AnalyticsHandler) GetRevenueTrend(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	start, end, _, err := parseAnalyticsParams(r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	points, err := h.svc.GetRevenueTrend(farmerID, start, end)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, points)
}

func (h *AnalyticsHandler) GetStockTrend(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	start, end, productID, err := parseAnalyticsParams(r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if productID == nil {
		utils.Error(w, http.StatusBadRequest, "product_id is required for stock trend")
		return
	}

	points, err := h.svc.GetStockTrend(farmerID, *productID, start, end)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, points)
}

func (h *AnalyticsHandler) GetRejectionTrend(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	start, end, productID, err := parseAnalyticsParams(r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if productID == nil {
		utils.Error(w, http.StatusBadRequest, "product_id is required for rejection trend")
		return
	}

	points, err := h.svc.GetRejectionTrend(farmerID, *productID, start, end)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, points)
}

func (h *AnalyticsHandler) GetProductSummary(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	start, end, _, err := parseAnalyticsParams(r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	summaries, err := h.svc.GetProductSummary(farmerID, start, end)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, summaries)
}

func (h *AnalyticsHandler) GetPlannedVsActual(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	start, end, productID, err := parseAnalyticsParams(r)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if productID == nil {
		utils.Error(w, http.StatusBadRequest, "product_id is required for planned vs actual")
		return
	}

	points, err := h.svc.GetPlannedVsActual(farmerID, *productID, start, end)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, points)
}

func (h *AnalyticsHandler) GetTodaySummary(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	summary, err := h.svc.GetTodaySummary(farmerID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, summary)
}