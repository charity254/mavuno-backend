package api

import (
	"encoding/json"
	"net/http"

	"github.com/mavuno/mavuno-backend/internal/services"
	"github.com/mavuno/mavuno-backend/internal/utils"
)

type MessageHandler struct {
	svc *services.MessageService
}

func NewMessageHandler(svc *services.MessageService) *MessageHandler {
	return &MessageHandler{svc: svc}
}

type sendMessageReq struct {
	Content string `json:"content"`
}

func (h *MessageHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	listingID, err := getUUIDParam(r, "id")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid listing ID")
		return
	}

	var req sendMessageReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	message, err := h.svc.SendMessage(listingID, userID, req.Content)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, message)
}

func (h *MessageHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	userID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	listingID, err := getUUIDParam(r, "id")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid listing ID")
		return
	}

	messages, err := h.svc.GetMessages(listingID, userID)
	if err != nil {
		utils.Error(w, http.StatusForbidden, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, messages)
}