package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"spectra-backend/internal/models"
	"spectra-backend/internal/repository"
)

// ConnectionHandler handles HTTP requests for connections (splicing)
type ConnectionHandler struct {
	repo *repository.ConnectionRepository
}

// NewConnectionHandler creates a new ConnectionHandler
func NewConnectionHandler(repo *repository.ConnectionRepository) *ConnectionHandler {
	return &ConnectionHandler{repo: repo}
}

// Create handles POST /api/connections
func (h *ConnectionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateConnectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Validate required fields
	if req.InputType == "" {
		respondError(w, http.StatusBadRequest, "InputType is required")
		return
	}
	if req.OutputType == "" {
		respondError(w, http.StatusBadRequest, "OutputType is required")
		return
	}
	if req.InputID == 0 {
		respondError(w, http.StatusBadRequest, "InputID is required")
		return
	}
	if req.OutputID == 0 {
		respondError(w, http.StatusBadRequest, "OutputID is required")
		return
	}

	connection, err := h.repo.Create(r.Context(), &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create connection: "+err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, models.SuccessResponse(connection, "Connection created successfully"))
}

// GetByID handles GET /api/connections/{id}
func (h *ConnectionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid connection ID")
		return
	}

	connection, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get connection: "+err.Error())
		return
	}

	if connection == nil {
		respondError(w, http.StatusNotFound, "Connection not found")
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse(connection, ""))
}

// List handles GET /api/connections
func (h *ConnectionHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := &models.ConnectionFilter{}

	// Parse query parameters
	if nodeIDParam := r.URL.Query().Get("location_node_id"); nodeIDParam != "" {
		if id, err := strconv.ParseInt(nodeIDParam, 10, 64); err == nil {
			filter.LocationNodeID = &id
		}
	}
	if inputTypeParam := r.URL.Query().Get("input_type"); inputTypeParam != "" {
		inputType := models.ConnectionType(inputTypeParam)
		filter.InputType = &inputType
	}
	if inputIDParam := r.URL.Query().Get("input_id"); inputIDParam != "" {
		if id, err := strconv.ParseInt(inputIDParam, 10, 64); err == nil {
			filter.InputID = &id
		}
	}
	if outputTypeParam := r.URL.Query().Get("output_type"); outputTypeParam != "" {
		outputType := models.ConnectionType(outputTypeParam)
		filter.OutputType = &outputType
	}
	if outputIDParam := r.URL.Query().Get("output_id"); outputIDParam != "" {
		if id, err := strconv.ParseInt(outputIDParam, 10, 64); err == nil {
			filter.OutputID = &id
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

	connections, total, err := h.repo.List(r.Context(), filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list connections: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.NewPaginatedResponse(connections, total, filter.Limit, filter.Offset))
}

// Delete handles DELETE /api/connections/{id}
func (h *ConnectionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid connection ID")
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		if err.Error() == "connection not found" {
			respondError(w, http.StatusNotFound, "Connection not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to delete connection: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse(nil, "Connection deleted successfully"))
}

// GetSpliceMatrix handles GET /api/connections/matrix/{nodeId}
func (h *ConnectionHandler) GetSpliceMatrix(w http.ResponseWriter, r *http.Request) {
	nodeID, err := getIDFromPath(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid node ID")
		return
	}

	matrix, err := h.repo.GetSpliceMatrix(r.Context(), nodeID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get splice matrix: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse(matrix, ""))
}

// GetByLocation handles GET /api/connections/location/{nodeId}
func (h *ConnectionHandler) GetByLocation(w http.ResponseWriter, r *http.Request) {
	nodeID, err := getIDFromPath(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid node ID")
		return
	}

	connections, err := h.repo.GetByLocation(r.Context(), nodeID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get connections by location: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse(connections, ""))
}
