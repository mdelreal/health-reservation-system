package storage

import (
	"time"

	"gorm.io/gorm"

	"github.com/manueldelreal/health-reservation-system/internal/models"
)

// AddAvailabilityAndSlots saves availability and corresponding slots to the database in a single transaction.
func AddAvailabilityAndSlots(providerID string, availabilities []models.Availability, slots []models.Slot) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// Save availabilities
		if err := tx.Create(&availabilities).Error; err != nil {
			return err
		}

		// Save slots
		if err := tx.Create(&slots).Error; err != nil {
			return err
		}

		return nil
	})
}

func GetAvailableSlots(providerID string, date time.Time) ([]models.Slot, error) {
	var slots []models.Slot

	start := date
	end := date.Add(24 * time.Hour)

	err := DB.Where("availability_id IN (?) AND status = ?",
		DB.Model(&models.Availability{}).Select("id").Where("provider_id = ? AND start_time BETWEEN ? AND ?", providerID, start, end),
		"Available",
	).Find(&slots).Error

	return slots, err
}

func ReserveSlot(slotID, clientID string, expiration time.Time) error {
	return DB.Model(&models.Slot{}).Where("id = ? AND status = ?", slotID, "Available").
		Updates(models.Slot{
			Status:            "Reserved",
			ReservationExpiry: &expiration,
		}).Error
}

func ConfirmReservation(slotID string) error {
	return DB.Model(&models.Slot{}).Where("id = ?", slotID).Updates(models.Slot{
		Status: "Confirmed",
	}).Error
}

func CleanupExpiredReservations() error {
	now := time.Now()
	return DB.Model(&models.Slot{}).Where("status = ? AND reservation_expiry < ?", "Reserved", now).
		Update("status", "Available").Error
}
