package api

import (
	"database/sql"
	"net/http"


	"github.com/gorilla/mux"
	"github.com/mavuno/mavuno-backend/internal/config"
    "github.com/mavuno/mavuno-backend/internal/utils"
)

func NewRouter(db *sql.DB, cfg *config.Config) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		utils.JSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	}).Methods("GET")
	return router
}