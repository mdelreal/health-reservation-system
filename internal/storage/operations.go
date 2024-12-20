package storage

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/manueldelreal/health-reservation-system/internal/models"
)

// AddAvailabilityAndSlots saves availability and corresponding slots to the database in a single transaction.
func AddAvailabilityAndSlots(providerID string, availabilities []models.Availability, slots []models.Slot) error {
	for i := range slots {
		slots[i].ProviderID = providerID
	}

	return DB.Transaction(func(tx *gorm.DB) error {
		txHandler := &GormDBHandler{DB: tx}

		// Save availabilities
		if err := txHandler.Create(&availabilities); err != nil {
			return err
		}

		// Save slots
		if err := txHandler.Create(&slots); err != nil {
			return err
		}

		return nil
	})
}

func GetAvailableSlots(providerID string, date time.Time) ([]models.Slot, error) {
	var slots []models.Slot

	start := date
	end := date.Add(24 * time.Hour)

	// Use the First method for the subquery and Find for the main query
	err := DB.(*GormDBHandler).GetDB().Where(
		"availability_id IN (?) AND status = ?",
		DB.(*GormDBHandler).GetDB().Model(&models.Availability{}).Select("id").Where(
			"provider_id = ? AND start_time BETWEEN ? AND ?", providerID, start, end,
		),
		"Available",
	).Find(&slots).Error

	return slots, err
}

func ReserveSlot(slotID, clientID, providerID string, expiration time.Time) error {
	return DB.(*GormDBHandler).GetDB().Transaction(func(tx *gorm.DB) error {
		// Fetch the slot to validate availability
		var slot models.Slot
		if err := tx.First(&slot, "id = ? AND status = ?", slotID, "Available").Error; err != nil {
			return err
		}

		// Create a reservation using the same slot ID
		reservation := models.Reservation{
			ID:                slot.ID, // Use the slot ID as the reservation ID
			ClientID:          clientID,
			ProviderID:        slot.ProviderID,
			AvailabilityID:    slot.AvailabilityID,
			ReservationExpiry: &expiration,
			StartTime:         slot.StartTime,
			EndTime:           slot.EndTime,
			Status:            "Reserved",
		}
		if err := tx.Create(&reservation).Error; err != nil {
			return err
		}

		// Remove the slot from the slots table
		if err := tx.Delete(&slot).Error; err != nil {
			return err
		}

		return nil
	})
}

func ConfirmReservation(reservationID string) error {
	// Fetch the reservation
	var reservation models.Reservation
	result := DB.(*GormDBHandler).GetDB().First(&reservation, "id = ?", reservationID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("reservation not found")
		}
		return result.Error
	}

	// Check if the reservation is already confirmed
	if reservation.Status == "Confirmed" {
		return errors.New("reservation is already confirmed")
	}

	// Update the reservation status to Confirmed
	return DB.(*GormDBHandler).GetDB().Model(&models.Reservation{}).Where("id = ?", reservationID).
		Update("status", "Confirmed").Error
}

func CleanupExpiredReservations() error {
	now := time.Now()
	log.Println("Checking for expired reservations...")

	return DB.(*GormDBHandler).GetDB().Transaction(func(tx *gorm.DB) error {
		// Fetch expired reservations
		var expiredReservations []models.Reservation
		if err := tx.Where("status = ? AND reservation_expiry < ?", "Reserved", now).Find(&expiredReservations).Error; err != nil {
			return err
		}

		// Iterate over expired reservations
		for _, reservation := range expiredReservations {
			// Recreate the slot with "Available" status in the slots table
			newSlot := models.Slot{
				ID:             reservation.ID, // Use the same ID
				AvailabilityID: reservation.AvailabilityID,
				ProviderID:     reservation.ProviderID,
				StartTime:      reservation.StartTime,
				EndTime:        reservation.EndTime,
				Status:         "Available",
			}

			if err := tx.Create(&newSlot).Error; err != nil {
				return err
			}

			// Delete the reservation
			if err := tx.Delete(&models.Reservation{}, "id = ?", reservation.ID).Error; err != nil {
				return err
			}
		}

		log.Printf("Expired reservations cleaned up: %d records processed", len(expiredReservations))
		return nil
	})
}

func GetReservationsByProvider(providerID string, date *time.Time) ([]models.Reservation, error) {
	var reservations []models.Reservation
	query := DB.(*GormDBHandler).GetDB().Where("provider_id = ?", providerID)

	if date != nil {
		start := date.Truncate(24 * time.Hour)
		end := start.Add(24 * time.Hour)
		query = query.Where("slot_id IN (?)", DB.(*GormDBHandler).GetDB().
			Model(&models.Slot{}).Select("id").Where("start_time BETWEEN ? AND ?", start, end))
	}

	err := query.Preload("Slot").Find(&reservations).Error
	return reservations, err
}

func GetReservationsByClient(clientID string, date *time.Time) ([]models.Reservation, error) {
	var reservations []models.Reservation
	query := DB.(*GormDBHandler).GetDB().Where("client_id = ?", clientID)

	if date != nil {
		start := date.Truncate(24 * time.Hour)
		end := start.Add(24 * time.Hour)
		query = query.Where("slot_id IN (?)", DB.(*GormDBHandler).GetDB().
			Model(&models.Slot{}).Select("id").Where("start_time BETWEEN ? AND ?", start, end))
	}

	err := query.Preload("Slot").Find(&reservations).Error
	return reservations, err
}
