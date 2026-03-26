package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mavuno/mavuno-backend/internal/middleware"
	"github.com/mavuno/mavuno-backend/internal/services"
	"github.com/mavuno/mavuno-backend/internal/utils"
)

type EntryHandler struct {
	svc *services.ProduceEntryService
}

func NewProduceEntryHandler(svc *services.ProduceEntryService) *EntryHandler {
	return &EntryHandler{svc: svc}
}

type createReq struct { //expected JSON body for creating an entry
	ProductID 			string 	`json:"product_id"`
	EntryDate 			string 	`json:"entry_date"` // format: YYYY-MM-DD
	OpeningStock 		int 	`json:"opening_stock"`
	AddedStock 			int 	`json:"added_stock"`
	SoldQuantity 		int 	`json:"sold_quantity"`
	RejectedQuantity 	int 	`json:"rejected_quantity"`
	PricePerUnit 		int 	`json:"price_per_unit"` //in KES cents
	Notes 				string 	`json:"notes"` 

}

type updateReq struct {
	ProductID 			string 	`json:"product_id"`
	EntryDate			string	`json:"entry_date"`
	OpeningStock		int		`json:"opening_stock"`
	AddedStock			int		`json:"added_stock"`
	SoldQuantity		int		`json:"sold_quantity"`
	RejectedQuantity	int		`json:"rejected_quantity"`
	PricePerUnit		int		`json:"price_per_unit"`
	Notes				string	`json:"notes"`
	Version 			int		`json:"version"`
}

func(h *EntryHandler) CreateEntry(w http.ResponseWriter, r *http.Request) { //handles POST /api/entries
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	var req createReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Parse product ID from string to UUID
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	// Parse entry date — frontend must send YYYY-MM-DD format
	date, err := time.Parse("2006-01-02", req.EntryDate)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid date format — use YYYY-MM-DD")
		return
	}

	entry, err := h.svc.CreateEntry(
		farmerID,
		productID,
		date,
		req.OpeningStock,
		req.AddedStock,
		req.SoldQuantity,
		req.RejectedQuantity,
		req.PricePerUnit,
		req.Notes,
	)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, entry)
}

// GetEntries handles GET /api/entries
// Supports optional query params: start, end, product_id
func (h *EntryHandler) GetEntries(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Read optional query params from URL
	var start, end time.Time
	startS := r.URL.Query().Get("start")
	endS := r.URL.Query().Get("end")

	if startS != "" {
		start, err = time.Parse("2006-01-02", startS)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "invalid start date — use YYYY-MM-DD")
			return
		}
	}

	if endS != "" {
		end, err = time.Parse("2006-01-02", endS)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "invalid end date — use YYYY-MM-DD")
			return
		}
	}

	// Parse optional product ID filter
	var productID *uuid.UUID
	pIDStr := r.URL.Query().Get("product_id")
	if pIDStr != "" {
		pid, err := uuid.Parse(pIDStr)
		if err != nil {
			utils.Error(w, http.StatusBadRequest, "invalid product ID")
			return
		}
		productID = &pid
	}

	entries, err := h.svc.GetEntries(farmerID, start, end, productID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, entries)
}

// GetEntry handles GET /api/entries/{id}
func (h *EntryHandler) GetEntry(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	eID, err := getUUIDParam(r, "id")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid entry ID")
		return
	}

	entry, err := h.svc.GetEntryByID(eID, farmerID)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, entry)
}

// UpdateEntry handles PUT /api/entries/{id}
func (h *EntryHandler) UpdateEntry(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	eID, err := getUUIDParam(r, "id")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid entry ID")
		return
	}

	var req updateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	date, err := time.Parse("2006-01-02", req.EntryDate)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid date format — use YYYY-MM-DD")
		return
	}

	entry, err := h.svc.UpdateEntry(
		eID,
		farmerID,
		productID,
		date,
		req.OpeningStock,
		req.AddedStock,
		req.SoldQuantity,
		req.RejectedQuantity,
		req.PricePerUnit,
		req.Version,
		req.Notes,
	)
	if err != nil {
		if err.Error() == "conflict: entry was updated by another session" {
			utils.Error(w, http.StatusConflict, err.Error())
			return
		}
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, entry)
}

// DeleteEntry handles DELETE /api/entries/{id}
func (h *EntryHandler) DeleteEntry(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	eID, err := getUUIDParam(r, "id")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid entry ID")
		return
	}

	if err := h.svc.DeleteEntry(eID, farmerID); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// getFarmerIDFromContext reads the user ID from the request context.
func getFarmerIDFromContext(r *http.Request) (uuid.UUID, error) {
	userIDStr, ok := r.Context().Value(middleware.ContextUserID).(string)
	if !ok || userIDStr == "" {
		return uuid.Nil, fmt.Errorf("unauthorized")
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID")
	}
	return userID, nil
}
