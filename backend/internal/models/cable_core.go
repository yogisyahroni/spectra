package models

import (
	"time"
)

// CoreStatus represents the status of a cable core
type CoreStatus string

const (
	CoreStatusVacant   CoreStatus = "VACANT"
	CoreStatusUsed     CoreStatus = "USED"
	CoreStatusReserved CoreStatus = "RESERVED"
	CoreStatusDamaged  CoreStatus = "DAMAGED"
)

// CableCore represents an individual fiber core within a cable
type CableCore struct {
	ID        int64      `json:"id" db:"id"`
	CableID   int64      `json:"cable_id" db:"cable_id"`
	CoreIndex int        `json:"core_index" db:"core_index"`
	TubeColor *string    `json:"tube_color,omitempty" db:"tube_color"`
	CoreColor *string    `json:"core_color,omitempty" db:"core_color"`
	Status    CoreStatus `json:"status" db:"status"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// CreateCableCoreRequest represents the request body for creating a cable core
type CreateCableCoreRequest struct {
	CableID   int64       `json:"cable_id" validate:"required"`
	CoreIndex int         `json:"core_index" validate:"required,min=1"`
	TubeColor *string     `json:"tube_color,omitempty"`
	CoreColor *string     `json:"core_color,omitempty"`
	Status    *CoreStatus `json:"status,omitempty"`
}

// UpdateCableCoreRequest represents the request body for updating a cable core
type UpdateCableCoreRequest struct {
	TubeColor *string     `json:"tube_color,omitempty"`
	CoreColor *string     `json:"core_color,omitempty"`
	Status    *CoreStatus `json:"status,omitempty" validate:"omitempty,oneof=VACANT USED RESERVED DAMAGED"`
}

// StandardTubeColors defines the standard 12-color tube sequence
var StandardTubeColors = []string{
	"Blue", "Orange", "Green", "Brown", "Slate", "White",
	"Red", "Black", "Yellow", "Violet", "Rose", "Aqua",
}

// StandardCoreColors defines the standard 12-color core sequence
var StandardCoreColors = []string{
	"Blue", "Orange", "Green", "Brown", "Slate", "White",
	"Red", "Black", "Yellow", "Violet", "Rose", "Aqua",
}

// GenerateCoresForCable creates core entries for a cable based on core count
func GenerateCoresForCable(cableID int64, coreCount int) []CableCore {
	cores := make([]CableCore, coreCount)
	tubeCount := len(StandardTubeColors)
	coreColorCount := len(StandardCoreColors)

	for i := 0; i < coreCount; i++ {
		tubeIndex := i / coreColorCount
		coreIndex := i % coreColorCount

		tubeColor := StandardTubeColors[tubeIndex%tubeCount]
		coreColor := StandardCoreColors[coreIndex]

		cores[i] = CableCore{
			CableID:   cableID,
			CoreIndex: i + 1,
			TubeColor: &tubeColor,
			CoreColor: &coreColor,
			Status:    CoreStatusVacant,
		}
	}

	return cores
}
