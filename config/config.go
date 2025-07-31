package config

import (
	"errors"
	"os"
)

// Config holds all configuration for the application
type Config struct {
	MongoURI        string
	MongoDatabase   string
	MongoCollection string
	OpenAIAPIKey    string
	Port            string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		MongoURI:        os.Getenv("MONGO_URI"),
		MongoDatabase:   os.Getenv("MONGO_DATABASE"),
		MongoCollection: os.Getenv("MONGO_COLLECTION"),
		OpenAIAPIKey:    os.Getenv("OPENAI_API_KEY"),
		Port:            os.Getenv("PORT"),
	}

	// Set defaults
	if config.MongoDatabase == "" {
		config.MongoDatabase = "moderation"
	}

	if config.MongoCollection == "" {
		config.MongoCollection = "results"
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	// Validate required fields
	if config.MongoURI == "" {
		return nil, errors.New("MONGO_URI environment variable is required")
	}

	if config.OpenAIAPIKey == "" {
		return nil, errors.New("OPENAI_API_KEY environment variable is required")
	}

	return config, nil
}
