package api

import (
	"database/sql"
	"net/http"


	"github.com/gorilla/mux"
	"github.com/mavuno/mavuno-backend/internal/config"
	"github.com/mavuno/mavuno-backend/internal/services"
    "github.com/mavuno/mavuno-backend/internal/storage"
    "github.com/mavuno/mavuno-backend/internal/utils"
)

func NewRouter(db *sql.DB, cfg *config.Config) *mux.Router {
	router := mux.NewRouter()

	// ── Health Check ────────────────────────────────────────
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.JSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	}).Methods("GET")

	// ── Auth Routes ─────────────────────────────────────────
	userRepo := storage.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := NewAuthHandler(authService)

	router.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST")
	return router
}