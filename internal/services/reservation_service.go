package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"

	pb "github.com/manueldelreal/health-reservation-system/api"
	"github.com/manueldelreal/health-reservation-system/internal/models"
	"github.com/manueldelreal/health-reservation-system/internal/storage"
)

type ReservationService struct{}

// generateID generates a new ULID as a string.
func generateID() string {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}

func (s *ReservationService) SetAvailability(ctx context.Context, req *pb.SetAvailabilityRequest) (*pb.SetAvailabilityResponse, error) {
	availabilityMap := make(map[string]models.Availability)
	var slots []models.Slot

	for _, timeSlot := range req.TimeSlots {
		startTime, err := time.Parse(time.RFC3339, timeSlot.StartTime)
		if err != nil {
			return nil, errors.New("invalid start time format")
		}
		endTime, err := time.Parse(time.RFC3339, timeSlot.EndTime)
		if err != nil {
			return nil, errors.New("invalid end time format")
		}

		key := fmt.Sprintf("%s-%s", startTime, endTime)
		if _, exists := availabilityMap[key]; exists {
			continue
		}

		availability := models.Availability{
			ID:         generateID(),
			ProviderID: req.ProviderId,
			StartTime:  startTime,
			EndTime:    endTime,
		}
		availabilityMap[key] = availability

		for t := startTime; t.Before(endTime); t = t.Add(15 * time.Minute) {
			slots = append(slots, models.Slot{
				ID:             generateID(),
				AvailabilityID: availability.ID,
				StartTime:      t,
				EndTime:        t.Add(15 * time.Minute),
				Status:         "Available",
			})
		}
	}

	availabilities := make([]models.Availability, 0, len(availabilityMap))
	for _, availability := range availabilityMap {
		availabilities = append(availabilities, availability)
	}

	err := storage.AddAvailabilityAndSlots(req.ProviderId, availabilities, slots)
	if err != nil {
		return nil, err
	}

	return &pb.SetAvailabilityResponse{Message: "Availability set successfully"}, nil
}

func (s *ReservationService) GetAvailableSlots(ctx context.Context, req *pb.GetAvailableSlotsRequest) (*pb.GetAvailableSlotsResponse, error) {
	// Parse the requested date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	// Query the database for slots
	slots, err := storage.GetAvailableSlots(req.ProviderId, date)
	if err != nil {
		return nil, err
	}

	// Convert database results to protobuf response
	var pbSlots []*pb.TimeSlot
	for _, slot := range slots {
		pbSlots = append(pbSlots, &pb.TimeSlot{
			Id:        slot.ID,
			StartTime: slot.StartTime.Format(time.RFC3339),
			EndTime:   slot.EndTime.Format(time.RFC3339),
			Status:    slot.Status,
		})
	}

	return &pb.GetAvailableSlotsResponse{Slots: pbSlots}, nil
}

func (s *ReservationService) ReserveSlot(ctx context.Context, req *pb.ReserveSlotRequest) (*pb.ReserveSlotResponse, error) {
	// Fetch the slot to validate
	var slot models.Slot
	err := storage.DB.First(&slot, "id = ? AND status = ?", req.SlotId, "Available").Error
	if err != nil {
		return nil, errors.New("slot is not available")
	}

	// Validate 24-hour rule
	if slot.StartTime.Before(time.Now().Add(24 * time.Hour)) {
		return nil, errors.New("reservations must be made at least 24 hours in advance")
	}

	// Reserve the slot
	expiration := time.Now().Add(30 * time.Minute)
	err = storage.ReserveSlot(req.SlotId, req.ClientId, expiration)
	if err != nil {
		return nil, err
	}

	return &pb.ReserveSlotResponse{
		ReservationId: req.SlotId,
		Message:       "Slot reserved successfully",
	}, nil
}

func (s *ReservationService) ConfirmReservation(ctx context.Context, req *pb.ConfirmReservationRequest) (*pb.ConfirmReservationResponse, error) {
	// Confirm the reservation in the database
	err := storage.ConfirmReservation(req.ReservationId)
	if err != nil {
		return nil, err
	}

	return &pb.ConfirmReservationResponse{Message: "Reservation confirmed"}, nil
}
