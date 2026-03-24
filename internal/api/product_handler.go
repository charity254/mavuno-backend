package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mavuno/mavuno-backend/internal/middleware"
	"github.com/mavuno/mavuno-backend/internal/services"
	"github.com/mavuno/mavuno-backend/internal/utils"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

type createProductRequest struct { //expected JSON body for creating a product
	Name 		string `json:"name"`
	UnitType 	string `json:"unit_type"`
	Description string `json:"description"`
}

type updateProductRequest struct {  //expected JSON body for updating a product
	Name 		string `json:"name"`
	UnitType 	string `json:"unit_type"`
	Description string `json:"description"`
	Version 	int    `json:"version"` //match current version in database
}
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) { //handles POST/api/products...only farmers can create products
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}
	products, err := h.productService.GetProductsByFarmer(farmerID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return 
	}
	utils.JSON(w, http.StatusOK, products)
}

func(h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) { //handles GET /api/products/{id}. Returns a singlr product by id
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

    // Extract the product ID from the URL 
	productID, err := getUUIDParam(r, "id")
    if err != nil {
        utils.Error(w, http.StatusBadRequest, "invalid product ID")
        return
    }
	
	product, err := h.productService.GetProductByID(productID, farmerID)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, product)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
	}
	productID, err := getUUIDParam(r, "id")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid product ID")
		return 
	}
	var req updateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	product, err := h.productService.UpdateProduct(productID, farmerID, req.Name, req.UnitType, req.Description, req.Version)
	if err != nil {
		if err.Error() == "conflict product was updated by another session" {
			utils.Error(w, http.StatusConflict, err.Error()) //409 conflict
			return
		}
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {  //handles DELETE /api/products/{id}
	farmerID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}
	productID, err := getUUIDParam(r, "id")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid product ID")
		return
	}

	err = h.productService.DeleteProduct(productID, farmerID)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent) //204 No content (successfully deleted)
}

func getFarmerID(r *http.Request) (uuid.UUID, error) {
	userIDstr, ok := r.Context().Value(middleware.ContextUserID).(string)
	if !ok || userIDstr == "" {
		return uuid.Nil, fmt.Errorf("unauthorized")
	}
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID")
	}
	return userID, nil
}

func getUUIDParam(r *http.Request, param string) (uuid.UUID, error) {
	vars := mux.Vars(r)
	idStr, ok := vars[param]
	if !ok || idStr == "" {
		return uuid.Nil, fmt.Errorf("missing parameter: %s", param)
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid UUID: %s", param)
	}
	return id, nil
}
