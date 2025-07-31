package api

import (
	"net/http"
	"strings"

	"github.com/go-moderation-api/api/v1"
	"github.com/go-moderation-api/config"
)

// Router handles routing to the appropriate API version
type Router struct {
	v1Handler *v1.Handler
}

// NewRouter creates a new API router
func NewRouter(cfg *config.Config) (*Router, error) {
	// Initialize v1 handler
	v1Handler, err := v1.NewHandler(cfg)
	if err != nil {
		return nil, err
	}

	return &Router{
		v1Handler: v1Handler,
	}, nil
}

// ServeHTTP implements the http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	// Route to health check
	if path == "/api/health" {
		r.v1Handler.HealthCheck(w, req)
		return
	}

	// Handle versioned endpoints
	if strings.HasPrefix(path, "/api/v1/") {
		// Remove version prefix
		req.URL.Path = strings.Replace(path, "/api/v1", "", 1)
		r.handleV1(w, req)
		return
	}

	// Default to v1 for unversioned endpoints (for backward compatibility)
	if strings.HasPrefix(path, "/api/") {
		r.handleV1(w, req)
		return
	}

	// Handle unknown paths
	http.NotFound(w, req)
}

// handleV1 routes the request to the v1 API handler
func (r *Router) handleV1(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	switch path {
	case "/api/moderate":
		r.v1Handler.ModerateContent(w, req)
	default:
		http.NotFound(w, req)
	}
}
