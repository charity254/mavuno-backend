package main

import (
	"fmt"
	"log"
	"net/http"

	
	"github.com/joho/godotenv"
	"github.com/mavuno/mavuno-backend/internal/api"
	"github.com/mavuno/mavuno-backend/internal/config"
	"github.com/mavuno/mavuno-backend/internal/middleware"
	"github.com/mavuno/mavuno-backend/internal/storage"
	"github.com/rs/cors"
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
	// Apply global middleware
	c := cors.New(cors.Options{
    AllowedOrigins: []string{"*"},
    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders: []string{"Authorization", "Content-Type"},
})

	// Logging — logs every request with method, path, status and duration
	// Rate limiting — limits each IP to 100 requests per minute
	handler := middleware.LoggingMiddleware(
		middleware.RateLimitMiddleware(
			c.Handler(router),
			),
	)

	// Limit request body size to 1MB to prevent large payload attacks
	limitedHandler := http.MaxBytesHandler(handler, 1<<20)

	fmt.Printf("🌾 Mavuno server starting on port %s...\n", cfg.Port)

	if err := http.ListenAndServe(":"+cfg.Port, limitedHandler); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}