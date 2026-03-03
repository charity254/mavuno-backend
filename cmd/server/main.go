package main

import (
	"fmt"
	"log"
	"net/http"

	
	"github.com/joho/godotenv"
	"github.com/mavuno/mavuno-backend/internal/api"
	"github.com/mavuno/mavuno-backend/internal/config"
	"github.com/mavuno/mavuno-backend/internal/storage"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found - reading from system environment")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	db, err := storage.InitDB(cfg.DBUrl)
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}
	defer db.Close()

	router := api.NewRouter(db, cfg)

	fmt.Printf("🌾 Mavuno server starting on port %s...\n", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}