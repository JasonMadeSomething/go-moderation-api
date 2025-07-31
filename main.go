package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-moderation-api/api"
	"github.com/go-moderation-api/config"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file in development
	if os.Getenv("GO_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: No .env file found")
		}
	}

	// Initialize configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize API router
	router, err := api.NewRouter(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize API router: %v", err)
	}

	// Set up routes with the router as the handler
	http.Handle("/api/", router)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
