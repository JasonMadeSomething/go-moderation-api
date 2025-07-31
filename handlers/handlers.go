package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-moderation-api/config"
	"github.com/go-moderation-api/models"
	"github.com/go-moderation-api/utils"
)

// Handler handles HTTP requests
type Handler struct {
	mongodb    *utils.MongoDB
	openaiAPI  *utils.OpenAIClient
}

// NewHandler creates a new Handler
func NewHandler(cfg *config.Config) (*Handler, error) {
	// Initialize MongoDB client
	mongodb, err := utils.NewMongoDB(cfg)
	if err != nil {
		return nil, err
	}

	// Initialize OpenAI client
	openaiAPI := utils.NewOpenAIClient(cfg.OpenAIAPIKey)

	return &Handler{
		mongodb:    mongodb,
		openaiAPI:  openaiAPI,
	}, nil
}

// HealthCheck handles health check requests
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// ModerateContent handles content moderation requests
func (h *Handler) ModerateContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var request models.ModerationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if request.Content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	if request.SourceSystem == "" {
		http.Error(w, "Source system is required", http.StatusBadRequest)
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	// Check cache for existing moderation result
	result, err := h.mongodb.FindModerationResult(ctx, request.Content)
	if err != nil {
		log.Printf("Error checking cache: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// If result found in cache, return it
	if result != nil {
		response := models.ModerationResponse{
			Allowed: result.Allowed,
			Message: getResponseMessage(result.Allowed),
		}

		w.Header().Set("Content-Type", "application/json")
		if !result.Allowed {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Call OpenAI moderation API
	moderationResponse, err := h.openaiAPI.CheckModeration(request.Content)
	if err != nil {
		log.Printf("Error calling OpenAI API: %v", err)
		http.Error(w, "Error checking content moderation", http.StatusInternalServerError)
		return
	}

	// Determine if content is allowed
	allowed := !moderationResponse.Results[0].Flagged

	// Save result to cache
	cacheResult := &models.ModerationResult{
		Content:      request.Content,
		SourceSystem: request.SourceSystem,
		Allowed:      allowed,
		OpenAIResult: *moderationResponse,
	}

	if err := h.mongodb.SaveModerationResult(ctx, cacheResult); err != nil {
		log.Printf("Error saving to cache: %v", err)
		// Continue even if caching fails
	}

	// Return response
	response := models.ModerationResponse{
		Allowed: allowed,
		Message: getResponseMessage(allowed),
	}

	w.Header().Set("Content-Type", "application/json")
	if !allowed {
		w.WriteHeader(http.StatusForbidden)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(response)
}

// getResponseMessage returns an appropriate message based on moderation result
func getResponseMessage(allowed bool) string {
	if allowed {
		return "Content allowed"
	}
	return "Content violates content policy"
}
