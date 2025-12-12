package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"spectra-backend/internal/models"
	"spectra-backend/internal/repository"
)

// NodeHandler handles HTTP requests for nodes
type NodeHandler struct {
	repo *repository.NodeRepository
}

// NewNodeHandler creates a new NodeHandler
func NewNodeHandler(repo *repository.NodeRepository) *NodeHandler {
	return &NodeHandler{repo: repo}
}

// Create handles POST /api/nodes
func (h *NodeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Validate required fields
	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "Name is required")
		return
	}
	if req.Type == "" {
		respondError(w, http.StatusBadRequest, "Type is required")
		return
	}
	if req.Latitude == 0 && req.Longitude == 0 {
		respondError(w, http.StatusBadRequest, "Latitude and Longitude are required")
		return
	}

	node, err := h.repo.Create(r.Context(), &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create node: "+err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, models.SuccessResponse(node, "Node created successfully"))
}

// GetByID handles GET /api/nodes/{id}
func (h *NodeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid node ID")
		return
	}

	node, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get node: "+err.Error())
		return
	}

	if node == nil {
		respondError(w, http.StatusNotFound, "Node not found")
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse(node, ""))
}

// List handles GET /api/nodes
func (h *NodeHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := &models.NodeFilter{}

	// Parse query parameters
	if typeParam := r.URL.Query().Get("type"); typeParam != "" {
		nodeType := models.NodeType(typeParam)
		filter.Type = &nodeType
	}
	if statusParam := r.URL.Query().Get("status"); statusParam != "" {
		status := models.NodeStatus(statusParam)
		filter.Status = &status
	}
	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if limit, err := strconv.Atoi(limitParam); err == nil {
			filter.Limit = limit
		}
	}
	if offsetParam := r.URL.Query().Get("offset"); offsetParam != "" {
		if offset, err := strconv.Atoi(offsetParam); err == nil {
			filter.Offset = offset
		}
	}

	nodes, total, err := h.repo.List(r.Context(), filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list nodes: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.NewPaginatedResponse(nodes, total, filter.Limit, filter.Offset))
}

// Update handles PUT /api/nodes/{id}
func (h *NodeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid node ID")
		return
	}

	var req models.UpdateNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	node, err := h.repo.Update(r.Context(), id, &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update node: "+err.Error())
		return
	}

	if node == nil {
		respondError(w, http.StatusNotFound, "Node not found")
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse(node, "Node updated successfully"))
}

// Delete handles DELETE /api/nodes/{id}
func (h *NodeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid node ID")
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		if err.Error() == "node not found" {
			respondError(w, http.StatusNotFound, "Node not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to delete node: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse(nil, "Node deleted successfully"))
}

// GetNearby handles GET /api/nodes/nearby
func (h *NodeHandler) GetNearby(w http.ResponseWriter, r *http.Request) {
	latParam := r.URL.Query().Get("lat")
	lngParam := r.URL.Query().Get("lng")
	radiusParam := r.URL.Query().Get("radius")

	if latParam == "" || lngParam == "" {
		respondError(w, http.StatusBadRequest, "lat and lng are required")
		return
	}

	lat, err := strconv.ParseFloat(latParam, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid latitude")
		return
	}

	lng, err := strconv.ParseFloat(lngParam, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid longitude")
		return
	}

	radius := 1.0 // Default 1km
	if radiusParam != "" {
		if r, err := strconv.ParseFloat(radiusParam, 64); err == nil {
			radius = r / 1000 // Convert meters to km
		}
	}

	query := &models.NearbyQuery{
		Latitude:  lat,
		Longitude: lng,
		RadiusKM:  radius,
	}

	if typeParam := r.URL.Query().Get("type"); typeParam != "" {
		nodeType := models.NodeType(typeParam)
		query.Type = &nodeType
	}

	nodes, err := h.repo.GetNearby(r.Context(), query)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get nearby nodes: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse(nodes, ""))
}

// GetGeoJSON handles GET /api/geojson/nodes
func (h *NodeHandler) GetGeoJSON(w http.ResponseWriter, r *http.Request) {
	filter := &models.NodeFilter{
		Limit: 10000, // High limit for GeoJSON export
	}

	if typeParam := r.URL.Query().Get("type"); typeParam != "" {
		nodeType := models.NodeType(typeParam)
		filter.Type = &nodeType
	}

	features, err := h.repo.GetAllAsGeoJSON(r.Context(), filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get nodes as GeoJSON: "+err.Error())
		return
	}

	// Convert to interface slice for FeatureCollection
	interfaceFeatures := make([]interface{}, len(features))
	for i, f := range features {
		interfaceFeatures[i] = f
	}

	fc := models.NewGeoJSONFeatureCollection(interfaceFeatures)
	respondJSON(w, http.StatusOK, fc)
}
