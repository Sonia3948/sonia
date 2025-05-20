
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	MongoURI    string
	Port        string
	FrontendURL string
}

// Load reads configuration from environment variables
func Load() *Config {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}
	
	// Configure MongoDB connection
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017/elsaidaliya"
	}
	
	// Configure server port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	// Configure frontend URL for CORS
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	
	return &Config{
		MongoURI:    mongoURI,
		Port:        port,
		FrontendURL: frontendURL,
	}
}
