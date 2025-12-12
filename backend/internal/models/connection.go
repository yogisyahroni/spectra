package models

import (
	"time"
)

// ConnectionType represents the type of connection endpoint
type ConnectionType string

const (
	ConnectionTypeCore ConnectionType = "CORE"
	ConnectionTypePort ConnectionType = "PORT"
)

// Connection represents a splicing link between cable cores or ports
type Connection struct {
	ID             int64          `json:"id" db:"id"`
	LocationNodeID *int64         `json:"location_node_id,omitempty" db:"location_node_id"`
	InputType      ConnectionType `json:"input_type" db:"input_type"`
	InputID        int64          `json:"input_id" db:"input_id"`
	OutputType     ConnectionType `json:"output_type" db:"output_type"`
	OutputID       int64          `json:"output_id" db:"output_id"`
	LossDB         *float64       `json:"loss_db,omitempty" db:"loss_db"`
	Notes          *string        `json:"notes,omitempty" db:"notes"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`

	// Joined data
	LocationNode *Node `json:"location_node,omitempty" db:"-"`
}

// CreateConnectionRequest represents the request body for creating a connection
type CreateConnectionRequest struct {
	LocationNodeID *int64         `json:"location_node_id,omitempty"`
	InputType      ConnectionType `json:"input_type" validate:"required,oneof=CORE PORT"`
	InputID        int64          `json:"input_id" validate:"required"`
	OutputType     ConnectionType `json:"output_type" validate:"required,oneof=CORE PORT"`
	OutputID       int64          `json:"output_id" validate:"required"`
	LossDB         *float64       `json:"loss_db,omitempty"`
	Notes          *string        `json:"notes,omitempty"`
}

// UpdateConnectionRequest represents the request body for updating a connection
type UpdateConnectionRequest struct {
	LocationNodeID *int64          `json:"location_node_id,omitempty"`
	InputType      *ConnectionType `json:"input_type,omitempty" validate:"omitempty,oneof=CORE PORT"`
	InputID        *int64          `json:"input_id,omitempty"`
	OutputType     *ConnectionType `json:"output_type,omitempty" validate:"omitempty,oneof=CORE PORT"`
	OutputID       *int64          `json:"output_id,omitempty"`
	LossDB         *float64        `json:"loss_db,omitempty"`
	Notes          *string         `json:"notes,omitempty"`
}

// ConnectionFilter represents query filters for listing connections
type ConnectionFilter struct {
	LocationNodeID *int64 `json:"location_node_id,omitempty"`
	InputType      *ConnectionType
	InputID        *int64
	OutputType     *ConnectionType
	OutputID       *int64
	Limit          int `json:"limit,omitempty"`
	Offset         int `json:"offset,omitempty"`
}

// SpliceMatrix represents a visual representation of connections at a location
type SpliceMatrix struct {
	LocationNode *Node              `json:"location_node"`
	InputCores   []CableCore        `json:"input_cores"`
	OutputCores  []CableCore        `json:"output_cores"`
	Connections  []MatrixConnection `json:"connections"`
}

// MatrixConnection represents a connection in the splice matrix
type MatrixConnection struct {
	InputCoreID  int64    `json:"input_core_id"`
	OutputCoreID int64    `json:"output_core_id"`
	LossDB       *float64 `json:"loss_db,omitempty"`
}
