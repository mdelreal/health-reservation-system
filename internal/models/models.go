package models

import "time"

type Provider struct {
	ID   string `gorm:"primaryKey"`
	Name string
}

type Availability struct {
	ID         string `gorm:"primaryKey"`
	ProviderID string `gorm:"index"`
	StartTime  time.Time
	EndTime    time.Time
	Provider   Provider `gorm:"foreignKey:ProviderID"`
}

// Specify the singular table name for Availability
func (Availability) TableName() string {
	return "availability"
}

type Slot struct {
	ID             string `gorm:"primaryKey"`
	AvailabilityID string `gorm:"index"`
	StartTime      time.Time
	EndTime        time.Time
	Status         string
	ReservationID  string
	Availability   Availability `gorm:"foreignKey:AvailabilityID"`
	ProviderID     string       `gorm:"index"`
}

type Reservation struct {
	ID                string `gorm:"primaryKey"`
	SlotID            string `gorm:"index"`
	ClientID          string `gorm:"index"` // Add an index for efficient querying
	ProviderID        string `gorm:"index"`
	Status            string // Pending, Confirmed
	ReservationExpiry *time.Time
	Slot              Slot `gorm:"foreignKey:SlotID"`
}
