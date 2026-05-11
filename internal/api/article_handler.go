package api

import (
	"encoding/json"
	"net/http"

	"github.com/mavuno/mavuno-backend/internal/services"
	"github.com/mavuno/mavuno-backend/internal/utils"
)

type ArticleHandler struct {
	svc *services.ArticleService
}

func NewArticleHandler(svc *services.ArticleService) *ArticleHandler {
	return &ArticleHandler{svc: svc}
}

type createArticleReq struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
}

func (h *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	authorID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	var req createArticleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	article, err := h.svc.CreateArticle(authorID, req.Title, req.Content, req.Category)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, article)
}

func (h *ArticleHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	// Optional category filter e.g. /api/articles?category=Disease Alerts
	category := r.URL.Query().Get("category")

	articles, err := h.svc.GetArticles(category)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, articles)
}

func (h *ArticleHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	articleID, err := getUUIDParam(r, "id")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid article ID")
		return
	}

	article, err := h.svc.GetArticleByID(articleID)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, article)
}

func (h *ArticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	authorID, err := getFarmerID(r)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	articleID, err := getUUIDParam(r, "id")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "invalid article ID")
		return
	}

	if err := h.svc.DeleteArticle(articleID, authorID); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}