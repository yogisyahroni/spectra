package models

import (
	"time"
)

// CableType represents the type of cable
type CableType string

const (
	CableTypeADSS CableType = "ADSS"
	CableTypeDuct CableType = "DUCT"
	CableTypeDrop CableType = "DROP"
)

// CableStatus represents the operational status of a cable
type CableStatus string

const (
	CableStatusActive      CableStatus = "ACTIVE"
	CableStatusMaintenance CableStatus = "MAINTENANCE"
	CableStatusPlan        CableStatus = "PLAN"
	CableStatusInactive    CableStatus = "INACTIVE"
)

// Cable represents a fiber optic cable connecting two nodes
type Cable struct {
	ID           int64       `json:"id" db:"id"`
	Name         *string     `json:"name,omitempty" db:"name"`
	Type         CableType   `json:"type" db:"type"`
	CoreCount    int         `json:"core_count" db:"core_count"`
	LengthMeter  *float64    `json:"length_meter,omitempty" db:"length_meter"`
	OriginNodeID *int64      `json:"origin_node_id,omitempty" db:"origin_node_id"`
	DestNodeID   *int64      `json:"dest_node_id,omitempty" db:"dest_node_id"`
	ColorHex     string      `json:"color_hex" db:"color_hex"`
	Status       CableStatus `json:"status" db:"status"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at" db:"updated_at"`

	// Populated for GeoJSON export
	PathCoordinates [][]float64 `json:"path_coordinates,omitempty" db:"-"`

	// Joined data
	OriginNode *Node `json:"origin_node,omitempty" db:"-"`
	DestNode   *Node `json:"dest_node,omitempty" db:"-"`
}

// CreateCableRequest represents the request body for creating a cable
type CreateCableRequest struct {
	Name            *string     `json:"name,omitempty"`
	Type            CableType   `json:"type" validate:"required,oneof=ADSS DUCT DROP"`
	CoreCount       int         `json:"core_count" validate:"required,min=1,max=288"`
	LengthMeter     *float64    `json:"length_meter,omitempty"`
	OriginNodeID    *int64      `json:"origin_node_id,omitempty"`
	DestNodeID      *int64      `json:"dest_node_id,omitempty"`
	PathCoordinates [][]float64 `json:"path_coordinates,omitempty"` // [[lng, lat], [lng, lat], ...]
	ColorHex        *string     `json:"color_hex,omitempty"`
	Status          CableStatus `json:"status,omitempty"`
}

// UpdateCableRequest represents the request body for updating a cable
type UpdateCableRequest struct {
	Name            *string      `json:"name,omitempty"`
	Type            *CableType   `json:"type,omitempty" validate:"omitempty,oneof=ADSS DUCT DROP"`
	CoreCount       *int         `json:"core_count,omitempty" validate:"omitempty,min=1,max=288"`
	LengthMeter     *float64     `json:"length_meter,omitempty"`
	OriginNodeID    *int64       `json:"origin_node_id,omitempty"`
	DestNodeID      *int64       `json:"dest_node_id,omitempty"`
	PathCoordinates [][]float64  `json:"path_coordinates,omitempty"`
	ColorHex        *string      `json:"color_hex,omitempty"`
	Status          *CableStatus `json:"status,omitempty" validate:"omitempty,oneof=ACTIVE MAINTENANCE PLAN INACTIVE"`
}

// CableFilter represents query filters for listing cables
type CableFilter struct {
	Type         *CableType   `json:"type,omitempty"`
	Status       *CableStatus `json:"status,omitempty"`
	OriginNodeID *int64       `json:"origin_node_id,omitempty"`
	DestNodeID   *int64       `json:"dest_node_id,omitempty"`
	Limit        int          `json:"limit,omitempty"`
	Offset       int          `json:"offset,omitempty"`
}

// CableGeoJSON represents a cable in GeoJSON format
type CableGeoJSON struct {
	Type       string                 `json:"type"`
	Geometry   GeoJSONLineString      `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

// GeoJSONLineString represents a GeoJSON LineString geometry
type GeoJSONLineString struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"` // [[lng, lat], [lng, lat], ...]
}

// ToGeoJSON converts a Cable to GeoJSON format
func (c *Cable) ToGeoJSON() CableGeoJSON {
	properties := map[string]interface{}{
		"id":         c.ID,
		"type":       c.Type,
		"core_count": c.CoreCount,
		"status":     c.Status,
		"color_hex":  c.ColorHex,
	}

	if c.Name != nil {
		properties["name"] = *c.Name
	}
	if c.LengthMeter != nil {
		properties["length_meter"] = *c.LengthMeter
	}
	if c.OriginNodeID != nil {
		properties["origin_node_id"] = *c.OriginNodeID
	}
	if c.DestNodeID != nil {
		properties["dest_node_id"] = *c.DestNodeID
	}

	coords := c.PathCoordinates
	if coords == nil {
		coords = [][]float64{}
	}

	return CableGeoJSON{
		Type: "Feature",
		Geometry: GeoJSONLineString{
			Type:        "LineString",
			Coordinates: coords,
		},
		Properties: properties,
	}
}
