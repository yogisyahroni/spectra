package models

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Total   int64 `json:"total"`
	Limit   int   `json:"limit"`
	Offset  int   `json:"offset"`
	HasMore bool  `json:"has_more"`
}

// GeoJSONFeatureCollection represents a GeoJSON FeatureCollection
type GeoJSONFeatureCollection struct {
	Type     string        `json:"type"`
	Features []interface{} `json:"features"`
}

// NewGeoJSONFeatureCollection creates a new GeoJSON FeatureCollection
func NewGeoJSONFeatureCollection(features []interface{}) GeoJSONFeatureCollection {
	if features == nil {
		features = []interface{}{}
	}
	return GeoJSONFeatureCollection{
		Type:     "FeatureCollection",
		Features: features,
	}
}

// SuccessResponse creates a success response
func SuccessResponse(data interface{}, message string) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse creates an error response
func ErrorResponse(err string) Response {
	return Response{
		Success: false,
		Error:   err,
	}
}

// NewPaginatedResponse creates a paginated response
func NewPaginatedResponse(data interface{}, total int64, limit, offset int) PaginatedResponse {
	hasMore := int64(offset+limit) < total
	return PaginatedResponse{
		Success: true,
		Data:    data,
		Pagination: Pagination{
			Total:   total,
			Limit:   limit,
			Offset:  offset,
			HasMore: hasMore,
		},
	}
}
