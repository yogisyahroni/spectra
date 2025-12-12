package models

import (
	"time"
)

// NodeType represents the type of network node
type NodeType string

const (
	NodeTypeOLT      NodeType = "OLT"
	NodeTypeODC      NodeType = "ODC"
	NodeTypeODP      NodeType = "ODP"
	NodeTypeClosure  NodeType = "CLOSURE"
	NodeTypePole     NodeType = "POLE"
	NodeTypeCustomer NodeType = "CUSTOMER"
)

// NodeStatus represents the operational status of a node
type NodeStatus string

const (
	NodeStatusActive      NodeStatus = "ACTIVE"
	NodeStatusMaintenance NodeStatus = "MAINTENANCE"
	NodeStatusPlan        NodeStatus = "PLAN"
	NodeStatusInactive    NodeStatus = "INACTIVE"
)

// Node represents a network infrastructure point (OLT, ODC, ODP, etc.)
type Node struct {
	ID            int64      `json:"id" db:"id"`
	Name          string     `json:"name" db:"name"`
	Type          NodeType   `json:"type" db:"type"`
	Latitude      float64    `json:"latitude" db:"latitude"`
	Longitude     float64    `json:"longitude" db:"longitude"`
	Address       *string    `json:"address,omitempty" db:"address"`
	CapacityPorts int        `json:"capacity_ports" db:"capacity_ports"`
	UsedPorts     int        `json:"used_ports" db:"used_ports"`
	Model         *string    `json:"model,omitempty" db:"model"`
	Status        NodeStatus `json:"status" db:"status"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// CreateNodeRequest represents the request body for creating a node
type CreateNodeRequest struct {
	Name          string     `json:"name" validate:"required,min=1,max=100"`
	Type          NodeType   `json:"type" validate:"required,oneof=OLT ODC ODP CLOSURE POLE CUSTOMER"`
	Latitude      float64    `json:"latitude" validate:"required,latitude"`
	Longitude     float64    `json:"longitude" validate:"required,longitude"`
	Address       *string    `json:"address,omitempty"`
	CapacityPorts *int       `json:"capacity_ports,omitempty"`
	Model         *string    `json:"model,omitempty"`
	Status        NodeStatus `json:"status,omitempty"`
}

// UpdateNodeRequest represents the request body for updating a node
type UpdateNodeRequest struct {
	Name          *string     `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Type          *NodeType   `json:"type,omitempty" validate:"omitempty,oneof=OLT ODC ODP CLOSURE POLE CUSTOMER"`
	Latitude      *float64    `json:"latitude,omitempty" validate:"omitempty,latitude"`
	Longitude     *float64    `json:"longitude,omitempty" validate:"omitempty,longitude"`
	Address       *string     `json:"address,omitempty"`
	CapacityPorts *int        `json:"capacity_ports,omitempty"`
	UsedPorts     *int        `json:"used_ports,omitempty"`
	Model         *string     `json:"model,omitempty"`
	Status        *NodeStatus `json:"status,omitempty" validate:"omitempty,oneof=ACTIVE MAINTENANCE PLAN INACTIVE"`
}

// NodeFilter represents query filters for listing nodes
type NodeFilter struct {
	Type   *NodeType   `json:"type,omitempty"`
	Status *NodeStatus `json:"status,omitempty"`
	Limit  int         `json:"limit,omitempty"`
	Offset int         `json:"offset,omitempty"`
}

// NearbyQuery represents parameters for searching nearby nodes
type NearbyQuery struct {
	Latitude  float64   `json:"latitude" validate:"required,latitude"`
	Longitude float64   `json:"longitude" validate:"required,longitude"`
	RadiusKM  float64   `json:"radius_km" validate:"required,min=0.1,max=100"`
	Type      *NodeType `json:"type,omitempty"`
}

// NodeGeoJSON represents a node in GeoJSON format
type NodeGeoJSON struct {
	Type       string                 `json:"type"`
	Geometry   GeoJSONPoint           `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

// GeoJSONPoint represents a GeoJSON point geometry
type GeoJSONPoint struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"` // [longitude, latitude]
}

// ToGeoJSON converts a Node to GeoJSON format
func (n *Node) ToGeoJSON() NodeGeoJSON {
	properties := map[string]interface{}{
		"id":             n.ID,
		"name":           n.Name,
		"type":           n.Type,
		"status":         n.Status,
		"capacity_ports": n.CapacityPorts,
		"used_ports":     n.UsedPorts,
	}

	if n.Address != nil {
		properties["address"] = *n.Address
	}
	if n.Model != nil {
		properties["model"] = *n.Model
	}

	return NodeGeoJSON{
		Type: "Feature",
		Geometry: GeoJSONPoint{
			Type:        "Point",
			Coordinates: []float64{n.Longitude, n.Latitude},
		},
		Properties: properties,
	}
}
