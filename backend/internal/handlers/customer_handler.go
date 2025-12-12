package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"spectra-backend/internal/models"
	"spectra-backend/internal/repository"
)

// CustomerHandler handles HTTP requests for customers
type CustomerHandler struct {
	repo *repository.CustomerRepository
}

// NewCustomerHandler creates a new CustomerHandler
func NewCustomerHandler(repo *repository.CustomerRepository) *CustomerHandler {
	return &CustomerHandler{repo: repo}
}

// Create handles POST /api/customers
func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Validate required fields
	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "Name is required")
		return
	}

	customer, err := h.repo.Create(r.Context(), &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create customer: "+err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, models.SuccessResponse(customer, "Customer created successfully"))
}

// GetByID handles GET /api/customers/{id}
func (h *CustomerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	customer, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get customer: "+err.Error())
		return
	}

	if customer == nil {
		respondError(w, http.StatusNotFound, "Customer not found")
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse(customer, ""))
}

// List handles GET /api/customers
func (h *CustomerHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := &models.CustomerFilter{}

	// Parse query parameters
	if nodeIDParam := r.URL.Query().Get("node_id"); nodeIDParam != "" {
		if id, err := strconv.ParseInt(nodeIDParam, 10, 64); err == nil {
			filter.NodeID = &id
		}
	}
	if statusParam := r.URL.Query().Get("status"); statusParam != "" {
		status := models.CustomerStatus(statusParam)
		filter.CurrentStatus = &status
	}
	if searchParam := r.URL.Query().Get("search"); searchParam != "" {
		filter.Search = &searchParam
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

	customers, total, err := h.repo.List(r.Context(), filter)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list customers: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.NewPaginatedResponse(customers, total, filter.Limit, filter.Offset))
}

// Update handles PUT /api/customers/{id}
func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	var req models.UpdateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	customer, err := h.repo.Update(r.Context(), id, &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update customer: "+err.Error())
		return
	}

	if customer == nil {
		respondError(w, http.StatusNotFound, "Customer not found")
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse(customer, "Customer updated successfully"))
}

// Delete handles DELETE /api/customers/{id}
func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		if err.Error() == "customer not found" {
			respondError(w, http.StatusNotFound, "Customer not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to delete customer: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse(nil, "Customer deleted successfully"))
}

// GetLOS handles GET /api/customers/los
func (h *CustomerHandler) GetLOS(w http.ResponseWriter, r *http.Request) {
	customers, err := h.repo.GetLOSCustomers(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to get LOS customers: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse(customers, ""))
}

// UpdateStatus handles PATCH /api/customers/{id}/status
func (h *CustomerHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPathAt(r, 2)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	var req struct {
		Status  models.CustomerStatus `json:"status"`
		RxPower *float64              `json:"rx_power,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if req.Status == "" {
		respondError(w, http.StatusBadRequest, "Status is required")
		return
	}

	if err := h.repo.UpdateStatus(r.Context(), id, req.Status, req.RxPower); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update customer status: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, models.SuccessResponse(nil, "Customer status updated successfully"))
}
