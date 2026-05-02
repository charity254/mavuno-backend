package api

import (
	"encoding/json"
	"net/http"

	"github.com/mavuno/mavuno-backend/internal/services"
	"github.com/mavuno/mavuno-backend/internal/utils"
)

type ListingHandler struct {
	svc *services.ListingService
}

func NewListingHandler(svc *services.ListingService) *ListingHandler {
	return &ListingHandler{svc: svc}
}

type createListingReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Quantity    int    `json:"quantity"`
	UnitType    string `json:"unit_type"`
	Location    string `json:"location"`
}

type updateListingReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Quantity    int    `json:"quantity"`
	UnitType    string `json:"unit_type"`
	Location    string `json:"location"`
	Version     int    `json:"version"`
}

func (h *ListingHandler) CreateListing(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	var req createListingReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	listing, err := h.svc.CreateListing(farmerID, req.Title, req.Description, req.Price, req.Quantity, req.UnitType, req.Location)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, listing)
}

func (h *ListingHandler) GetListings(w http.ResponseWriter, r *http.Request) {
	listings, err := h.svc.GetAllListings()
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, listings)
}

func (h *ListingHandler) GetListing(w http.ResponseWriter, r *http.Request) {
	listingID, err := getUUIDParam(r, "id")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid listing ID")
		return
	}

	listing, err := h.svc.GetListingByID(listingID)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, listing)
}

func (h *ListingHandler) UpdateListing(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	listingID, err := getUUIDParam(r, "id")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid listing ID")
		return
	}

	var req updateListingReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	listing, err := h.svc.UpdateListing(listingID, farmerID, req.Title, req.Description, req.Price, req.Quantity, req.UnitType, req.Location, req.Version)
	if err != nil {
		if err.Error() == "conflict: listing was updated by another session" {
			utils.Error(w, http.StatusConflict, err.Error())
			return
		}
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, listing)
}

func (h *ListingHandler) DeleteListing(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	listingID, err := getUUIDParam(r, "id")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid listing ID")
		return
	}

	if err := h.svc.DeleteListing(listingID, farmerID); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}