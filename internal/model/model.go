// Package model contains all entities using by trading service app
package model

import (
	"sync"

	"github.com/google/uuid"
)

// Share is a struct for shares entity
type Share struct {
	Name  string
	Price float64
}

// Position contains all info abount position, its name its ID, ID of profile opened this position vector, money and its borders
type Position struct {
	MoneyAmount float64
	ProfileID   uuid.UUID
	PositionID  uuid.UUID
	ShareName   string
	StopLoss    float64
	TakeProfit  float64
	Vector      string
}

// ProfilesManager is an assistent for managing all opened positions by profiles
type ProfilesManager struct {
	Mu sync.RWMutex

	// Profiles contains ProfileID as a key and slice of selected shares as value
	Profiles map[uuid.UUID][]string

	// ProfilesShares contains ProfileID as a key and channel with slice of shares to deliver it
	ProfilesShares map[uuid.UUID]chan []*Share

	// SharesPositionsRead contains map with ProfileID as a key and map of share and its count as value, this structs see how many positions read exact share
	SharesPositionsRead map[uuid.UUID]map[string]int

	// PositionClose contains map with PositionID as a value and bool channel as a value, when position closing field by PositionID gets true value in the chan
	PositionClose map[uuid.UUID]chan bool
}
