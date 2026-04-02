package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/mavuno/mavuno-backend/internal/services"
	"github.com/mavuno/mavuno-backend/internal/utils"
)

type SupplyAgreementHandler struct {
	svc *services.SupplyAgreementService
}

func NewSupplyAgreementHandler(svc *services.SupplyAgreementService) *SupplyAgreementHandler {
	return &SupplyAgreementHandler{svc: svc}
}

type createAgreementReq struct {
	ProductID        string   `json:"product_id"`
	LocationID       string   `json:"supply_location_id"`
	QtyPerDelivery   int      `json:"quantity_per_delivery"`
	PricePerUnit     int      `json:"price_per_unit"`
	DeliveryDays     []string `json:"delivery_days"`
	DeliveryNotes    string   `json:"delivery_notes"`
}

type updateAgreementReq struct {
	ProductID        string   `json:"product_id"`
	LocationID       string   `json:"supply_location_id"`
	QtyPerDelivery   int      `json:"quantity_per_delivery"`
	PricePerUnit     int      `json:"price_per_unit"`
	DeliveryDays     []string `json:"delivery_days"`
	DeliveryNotes    string   `json:"delivery_notes"`
	Active           bool     `json:"active"`
	Version          int      `json:"version"`
}

//handles POST /api/supply-agreements
func (h *SupplyAgreementHandler) CreateSupplyAgreement(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	var req createAgreementReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Parse product ID
	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	// Parse location ID
	locationID, err := uuid.Parse(req.LocationID)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid supply location ID")
		return
	}

	sa, err := h.svc.CreateSupplyAgreement(farmerID, productID, locationID, req.QtyPerDelivery, req.PricePerUnit, req.DeliveryDays, req.DeliveryNotes)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, sa)
}

// handles GET /api/supply-agreements
func (h *SupplyAgreementHandler) GetSupplyAgreements(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	agreements, err := h.svc.GetSupplyAgreements(farmerID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, agreements)
}

// handles GET /api/supply-agreements/active
// Returns only active agreements — used by frontend for dashboard reminders.
func (h *SupplyAgreementHandler) GetActiveSupplyAgreements(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	agreements, err := h.svc.GetActiveSupplyAgreements(farmerID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, agreements)
}

//handles GET /api/supply-agreements/{id}
func (h *SupplyAgreementHandler) GetSupplyAgreement(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	saID, err := getUUIDParam(r, "id")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid supply agreement ID")
		return
	}

	sa, err := h.svc.GetSupplyAgreementByID(saID, farmerID)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, sa)
}

//handles PUT /api/supply-agreements/{id}
func (h *SupplyAgreementHandler) UpdateSupplyAgreement(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	saID, err := getUUIDParam(r, "id")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid supply agreement ID")
		return
	}

	var req updateAgreementReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	locationID, err := uuid.Parse(req.LocationID)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid supply location ID")
		return
	}

	sa, err := h.svc.UpdateSupplyAgreement(saID, farmerID, productID, locationID, req.QtyPerDelivery, req.PricePerUnit, req.DeliveryDays, req.DeliveryNotes, req.Active, req.Version)
	if err != nil {
		if err.Error() == "conflict: supply agreement was updated by another session" {
			utils.Error(w, http.StatusConflict, err.Error())
			return
		}
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, sa)
}

//handles DELETE /api/supply-agreements/{id}
func (h *SupplyAgreementHandler) DeleteSupplyAgreement(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	saID, err := getUUIDParam(r, "id")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid supply agreement ID")
		return
	}

	if err := h.svc.DeleteSupplyAgreement(saID, farmerID); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}