package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ModerationRequest represents the incoming request for content moderation
type ModerationRequest struct {
	SourceSystem string `json:"source_system" bson:"source_system"`
	Content      string `json:"content" bson:"content"`
}

// OpenAIModerationRequest represents the request sent to OpenAI's moderation API
type OpenAIModerationRequest struct {
	Input string `json:"input"`
}

// OpenAIModerationResponse represents the response from OpenAI's moderation API
type OpenAIModerationResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Results []struct {
		Flagged        bool   `json:"flagged"`
		Categories     map[string]bool   `json:"categories"`
		CategoryScores map[string]float64 `json:"category_scores"`
	} `json:"results"`
}

// ModerationResult represents a cached moderation result in MongoDB
type ModerationResult struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Content      string             `bson:"content"`
	SourceSystem string             `bson:"source_system"`
	Allowed      bool               `bson:"allowed"`
	OpenAIResult OpenAIModerationResponse `bson:"openai_result,omitempty"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}

// ModerationResponse represents the API response for a moderation request
type ModerationResponse struct {
	Allowed bool   `json:"allowed"`
	Message string `json:"message,omitempty"`
	Version string `json:"version,omitempty"`
}
