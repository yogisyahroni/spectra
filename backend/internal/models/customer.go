package models

import (
	"time"
)

// CustomerStatus represents the connection status of a customer
type CustomerStatus string

const (
	CustomerStatusOnline   CustomerStatus = "ONLINE"
	CustomerStatusOffline  CustomerStatus = "OFFLINE"
	CustomerStatusLOS      CustomerStatus = "LOS"
	CustomerStatusPowerOff CustomerStatus = "POWER_OFF"
)

// Customer represents a fiber optic subscriber
type Customer struct {
	ID               int64          `json:"id" db:"id"`
	NodeID           *int64         `json:"node_id,omitempty" db:"node_id"`
	Name             string         `json:"name" db:"name"`
	ONTSN            *string        `json:"ont_sn,omitempty" db:"ont_sn"`
	Phone            *string        `json:"phone,omitempty" db:"phone"`
	Email            *string        `json:"email,omitempty" db:"email"`
	CurrentStatus    CustomerStatus `json:"current_status" db:"current_status"`
	LastRxPower      *float64       `json:"last_rx_power,omitempty" db:"last_rx_power"`
	SubscriptionType *string        `json:"subscription_type,omitempty" db:"subscription_type"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`

	// Joined data
	Node *Node `json:"node,omitempty" db:"-"`
}

// CreateCustomerRequest represents the request body for creating a customer
type CreateCustomerRequest struct {
	NodeID           *int64          `json:"node_id,omitempty"`
	Name             string          `json:"name" validate:"required,min=1,max=100"`
	ONTSN            *string         `json:"ont_sn,omitempty"`
	Phone            *string         `json:"phone,omitempty"`
	Email            *string         `json:"email,omitempty" validate:"omitempty,email"`
	CurrentStatus    *CustomerStatus `json:"current_status,omitempty"`
	SubscriptionType *string         `json:"subscription_type,omitempty"`
}

// UpdateCustomerRequest represents the request body for updating a customer
type UpdateCustomerRequest struct {
	NodeID           *int64          `json:"node_id,omitempty"`
	Name             *string         `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	ONTSN            *string         `json:"ont_sn,omitempty"`
	Phone            *string         `json:"phone,omitempty"`
	Email            *string         `json:"email,omitempty" validate:"omitempty,email"`
	CurrentStatus    *CustomerStatus `json:"current_status,omitempty" validate:"omitempty,oneof=ONLINE OFFLINE LOS POWER_OFF"`
	LastRxPower      *float64        `json:"last_rx_power,omitempty"`
	SubscriptionType *string         `json:"subscription_type,omitempty"`
}

// CustomerFilter represents query filters for listing customers
type CustomerFilter struct {
	NodeID        *int64          `json:"node_id,omitempty"`
	CurrentStatus *CustomerStatus `json:"current_status,omitempty"`
	Search        *string         `json:"search,omitempty"` // Search by name or ONT SN
	Limit         int             `json:"limit,omitempty"`
	Offset        int             `json:"offset,omitempty"`
}

// CustomerWithTrace represents a customer with their connection trace
type CustomerWithTrace struct {
	Customer   Customer    `json:"customer"`
	TracePath  []TraceNode `json:"trace_path"`
	TotalLoss  float64     `json:"total_loss_db"`
	TotalHops  int         `json:"total_hops"`
	TraceValid bool        `json:"trace_valid"`
}

// TraceNode represents a node in the connection trace
type TraceNode struct {
	Node     Node       `json:"node"`
	Cable    *Cable     `json:"cable,omitempty"`
	Core     *CableCore `json:"core,omitempty"`
	LossDB   float64    `json:"loss_db"`
	Sequence int        `json:"sequence"`
}

// RxPowerThreshold defines the acceptable Rx power levels
const (
	RxPowerGood     = -25.0 // dBm - Good signal
	RxPowerWarning  = -27.0 // dBm - Warning threshold
	RxPowerCritical = -28.0 // dBm - Critical (LOS imminent)
)

// GetRxPowerStatus returns the status based on Rx power level
func GetRxPowerStatus(rxPower float64) string {
	if rxPower >= RxPowerGood {
		return "GOOD"
	} else if rxPower >= RxPowerWarning {
		return "WARNING"
	}
	return "CRITICAL"
}
