package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBUrl string
	JWTSecret string
	Port string
}

func LoadConfig()(*Config, error) {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		return nil, fmt.Errorf("DB_URL environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		DBUrl: 		dbUrl,
		JWTSecret: 	jwtSecret,
		Port: 		port,
	}, nil
}