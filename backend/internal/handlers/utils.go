package handlers

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"strings"

	"spectra-backend/internal/models"
)

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// respondError sends an error response
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, models.ErrorResponse(message))
}

// getIDFromPath extracts the ID from the URL path
// Expects paths like /api/nodes/123
func getIDFromPath(r *http.Request) (int64, error) {
	// Get the last segment of the path
	urlPath := strings.TrimSuffix(r.URL.Path, "/")
	segment := path.Base(urlPath)
	return strconv.ParseInt(segment, 10, 64)
}

// getIDFromPathAt extracts an ID from a specific position in the URL path
// Useful for nested resources like /api/cables/123/cores/456
func getIDFromPathAt(r *http.Request, position int) (int64, error) {
	urlPath := strings.TrimPrefix(r.URL.Path, "/")
	urlPath = strings.TrimSuffix(urlPath, "/")
	segments := strings.Split(urlPath, "/")

	if position < 0 || position >= len(segments) {
		return 0, strconv.ErrSyntax
	}

	return strconv.ParseInt(segments[position], 10, 64)
}

// parseIntParam parses an integer query parameter with a default value
func parseIntParam(r *http.Request, key string, defaultValue int) int {
	param := r.URL.Query().Get(key)
	if param == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(param)
	if err != nil {
		return defaultValue
	}
	return value
}

// parseFloatParam parses a float query parameter with a default value
func parseFloatParam(r *http.Request, key string, defaultValue float64) float64 {
	param := r.URL.Query().Get(key)
	if param == "" {
		return defaultValue
	}
	value, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

// parseBoolParam parses a boolean query parameter with a default value
func parseBoolParam(r *http.Request, key string, defaultValue bool) bool {
	param := r.URL.Query().Get(key)
	if param == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(param)
	if err != nil {
		return defaultValue
	}
	return value
}
