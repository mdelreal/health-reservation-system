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

type Slot struct {
	ID                string `gorm:"primaryKey"`
	AvailabilityID    string `gorm:"index"`
	StartTime         time.Time
	EndTime           time.Time
	Status            string // Available, Reserved, Confirmed
	ReservationID     string `gorm:"index"`
	ReservationExpiry *time.Time
	Availability      Availability `gorm:"foreignKey:AvailabilityID"`
}

type Reservation struct {
	ID       string `gorm:"primaryKey"`
	SlotID   string `gorm:"index"`
	ClientID string
	Status   string // Pending, Confirmed
	Slot     Slot   `gorm:"foreignKey:SlotID"`
}
