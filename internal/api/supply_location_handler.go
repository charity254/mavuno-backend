package api

import (
	"encoding/json"
	"net/http"

	"github.com/mavuno/mavuno-backend/internal/services"
	"github.com/mavuno/mavuno-backend/internal/utils"
)

type SupplyLocationHandler struct {
	svc *services.SupplyLocationService
}

func NewSupplyLocationHandler(svc *services.SupplyLocationService) *SupplyLocationHandler {
	return &SupplyLocationHandler{svc: svc}
}

type createLocReq struct {
	Name            string `json:"name"`
	ContactPerson   string `json:"contact_person"`
	PhoneNumber     string `json:"phone_number"`
	LocationAddress string `json:"location_address"`
	Notes           string `json:"notes"`
}

type updateLocReq struct {
	Name            string `json:"name"`
	ContactPerson   string `json:"contact_person"`
	PhoneNumber     string `json:"phone_number"`
	LocationAddress string `json:"location_address"`
	Notes           string `json:"notes"`
	Version         int    `json:"version"`
}

func (h *SupplyLocationHandler) CreateSupplyLocation(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	var req createLocReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	loc, err := h.svc.CreateSupplyLocation(farmerID, req.Name, req.ContactPerson, req.PhoneNumber, req.LocationAddress, req.Notes)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, loc)
}

func (h *SupplyLocationHandler) GetSupplyLocations(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	locs, err := h.svc.GetSupplyLocations(farmerID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, locs)
}

func (h *SupplyLocationHandler) GetSupplyLocation(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	locID, err := getUUIDParam(r, "id")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid supply location ID")
		return
	}

	loc, err := h.svc.GetSupplyLocationByID(locID, farmerID)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, loc)
}

func (h *SupplyLocationHandler) UpdateSupplyLocation(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	locID, err := getUUIDParam(r, "id")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid supply location ID")
		return
	}

	var req updateLocReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	loc, err := h.svc.UpdateSupplyLocation(locID, farmerID, req.Name, req.ContactPerson, req.PhoneNumber, req.LocationAddress, req.Notes, req.Version)
	if err != nil {
		if err.Error() == "conflict: supply location was updated by another session" {
			utils.Error(w, http.StatusConflict, err.Error())
			return
		}
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, loc)
}

func (h *SupplyLocationHandler) DeleteSupplyLocation(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	locID, err := getUUIDParam(r, "id")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid supply location ID")
		return
	}

	if err := h.svc.DeleteSupplyLocation(locID, farmerID); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}