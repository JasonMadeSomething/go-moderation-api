package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-moderation-api/models"
)

const (
	openAIModerationEndpoint = "https://api.openai.com/v1/moderations"
	requestTimeout           = 10 * time.Second
)

// OpenAIClient represents an OpenAI API client
type OpenAIClient struct {
	apiKey string
	client *http.Client
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: requestTimeout,
		},
	}
}

// CheckModeration sends a request to the OpenAI moderation API
func (c *OpenAIClient) CheckModeration(content string) (*models.OpenAIModerationResponse, error) {
	requestBody := models.OpenAIModerationRequest{
		Input: content,
	}
	
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	
	req, err := http.NewRequest("POST", openAIModerationEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API returned non-200 status code: %d", resp.StatusCode)
	}
	
	var moderationResponse models.OpenAIModerationResponse
	if err := json.NewDecoder(resp.Body).Decode(&moderationResponse); err != nil {
		return nil, err
	}
	
	if len(moderationResponse.Results) == 0 {
		return nil, errors.New("OpenAI API returned empty results")
	}
	
	return &moderationResponse, nil
}
