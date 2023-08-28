// Package model contains all entities using by trading service app
package model

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// Share is a struct for shares entity
type Share struct {
	Price float64
	Name  string
}

// Position contains all info abount position, its name its ID, ID of profile opened this position vector, money and its borders
type Position struct {
	MoneyAmount float64 `validate:"gt=0"`
	StopLoss    float64 `validate:"gt=0"`
	TakeProfit  float64 `validate:"gt=0"`
	ProfileID   uuid.UUID
	PositionID  uuid.UUID
	ShareName   string `validate:"required,min=1"`
	ShortOrLong string `validate:"oneof=long short"`
	Close       chan bool
}

// OpenedPosition contains position struct and an additional fields that were created while transaction usually opens
type OpenedPosition struct {
	ShareStartPrice float64
	ShareEndPrice   float64
	ShareAmount     float64
	OpenedTime      time.Time
	Position
}

// ProfilesManager is an assistent for managing all opened positions by profiles
type ProfilesManager struct {
	Mu       sync.RWMutex
	Shares   map[string]float64
	Profiles map[uuid.UUID]map[uuid.UUID]OpenedPosition
}
