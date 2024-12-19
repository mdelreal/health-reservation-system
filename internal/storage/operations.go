package storage

import (
	"time"

	"github.com/manueldelreal/health-reservation-system/internal/models"
)

func AddAvailability(providerID string, slots []models.Availability) error {
	for _, slot := range slots {
		slot.ProviderID = providerID
		err := DB.Create(&slot).Error
		if err != nil {
			return err
		}
	}
	return nil
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
