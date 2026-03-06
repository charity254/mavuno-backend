package api

import (
	"encoding/json"
	"net/http"

	 "github.com/mavuno/mavuno-backend/internal/services"
    "github.com/mavuno/mavuno-backend/internal/utils"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authsService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authsService}
}

type registerRequest struct {
	Email 	 string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Role	 string `json:"role"`
}

type loginRequest struct {
	Email	  string `json:"email"`
	Password  string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	err := h.authService.Register(req.Email, req.Password, req.FullName, req.Role)
	if err != nil {
		if err.Error() == "email already exists" {
			utils.Error(w, http.StatusConflict, err.Error())
			return
		}
		utils.Error(w, http.StatusBadRequest, err.Error())
		return 
	}
	utils.JSON(w, http.StatusCreated, map[string]string{
		"message":"ACcount created successfully!",
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return 
	}

	token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return 
	}
	utils.JSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}