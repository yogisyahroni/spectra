package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"spectra-backend/internal/models"
	"spectra-backend/internal/repository"
)

// CableHandler handles HTTP requests for cables
type CableHandler struct {
	repo *repository.CableRepository
}

// NewCableHandler creates a new CableHandler
func NewCableHandler(repo *repository.CableRepository) *CableHandler {
	return &CableHandler{repo: repo}
}

// Create handles POST /api/cables
func (h *CableHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Validate required fields
	if req.Type == "" {
		respondError(w, http.StatusBadRequest, "Type is required")
		return
	}
	if req.CoreCount <= 0 {
		respondError(w, http.StatusBadRequest, "CoreCount must be greater than 0")
		return
	}

	cable, err := h.repo.Create(r.Context(), &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create cable: "+err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, models.SuccessResponse(cable, "Cable created successfully"))
}

// GetByID handles GET /api/cables/{id}
func (h *CableHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid cable ID")
		return
	}

	cable, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get cable: "+err.Error())
		return
	}

	if cable == nil {
		respondError(w, http.StatusNotFound, "Cable not found")
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse(cable, ""))
}

// List handles GET /api/cables
func (h *CableHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := &models.CableFilter{}

	// Parse query parameters
	if typeParam := r.URL.Query().Get("type"); typeParam != "" {
		cableType := models.CableType(typeParam)
		filter.Type = &cableType
	}
	if statusParam := r.URL.Query().Get("status"); statusParam != "" {
		status := models.CableStatus(statusParam)
		filter.Status = &status
	}
	if originParam := r.URL.Query().Get("origin_node_id"); originParam != "" {
		if id, err := strconv.ParseInt(originParam, 10, 64); err == nil {
			filter.OriginNodeID = &id
		}
	}
	if destParam := r.URL.Query().Get("dest_node_id"); destParam != "" {
		if id, err := strconv.ParseInt(destParam, 10, 64); err == nil {
			filter.DestNodeID = &id
		}
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

	cables, total, err := h.repo.List(r.Context(), filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list cables: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.NewPaginatedResponse(cables, total, filter.Limit, filter.Offset))
}

// Update handles PUT /api/cables/{id}
func (h *CableHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid cable ID")
		return
	}

	var req models.UpdateCableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	cable, err := h.repo.Update(r.Context(), id, &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update cable: "+err.Error())
		return
	}

	if cable == nil {
		respondError(w, http.StatusNotFound, "Cable not found")
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse(cable, "Cable updated successfully"))
}

// Delete handles DELETE /api/cables/{id}
func (h *CableHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid cable ID")
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		if err.Error() == "cable not found" {
			respondError(w, http.StatusNotFound, "Cable not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to delete cable: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse(nil, "Cable deleted successfully"))
}

// GetCores handles GET /api/cables/{id}/cores
func (h *CableHandler) GetCores(w http.ResponseWriter, r *http.Request) {
	// Extract cable ID from path (position 2 in /api/cables/{id}/cores)
	id, err := getIDFromPathAt(r, 2)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid cable ID")
		return
	}

	cores, err := h.repo.GetCores(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get cable cores: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse(cores, ""))
}

// UpdateCore handles PUT /api/cables/{id}/cores/{coreId}
func (h *CableHandler) UpdateCore(w http.ResponseWriter, r *http.Request) {
	// Extract IDs from path
	cableID, err := getIDFromPathAt(r, 2)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid cable ID")
		return
	}

	coreID, err := getIDFromPathAt(r, 4)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid core ID")
		return
	}

	var req models.UpdateCableCoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	core, err := h.repo.UpdateCore(r.Context(), cableID, coreID, &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update cable core: "+err.Error())
		return
	}

	if core == nil {
		respondError(w, http.StatusNotFound, "Cable core not found")
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse(core, "Cable core updated successfully"))
}

// GetGeoJSON handles GET /api/geojson/cables
func (h *CableHandler) GetGeoJSON(w http.ResponseWriter, r *http.Request) {
	features, err := h.repo.GetAllAsGeoJSON(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get cables as GeoJSON: "+err.Error())
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
