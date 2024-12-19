package services

import (
	"context"
	"errors"
	"time"

	pb "github.com/manueldelreal/health-reservation-system/api"
	"github.com/manueldelreal/health-reservation-system/internal/models"
	"github.com/manueldelreal/health-reservation-system/internal/storage"
)

type ReservationService struct{}

func (s *ReservationService) SetAvailability(ctx context.Context, req *pb.SetAvailabilityRequest) (*pb.SetAvailabilityResponse, error) {
	// Convert protobuf request to database models
	var availabilities []models.Availability
	for _, slot := range req.TimeSlots {
		startTime, err := time.Parse(time.RFC3339, slot.StartTime)
		if err != nil {
			return nil, errors.New("invalid start time format")
		}
		endTime, err := time.Parse(time.RFC3339, slot.EndTime)
		if err != nil {
			return nil, errors.New("invalid end time format")
		}
		availabilities = append(availabilities, models.Availability{
			ID:         slot.Id,
			ProviderID: req.ProviderId,
			StartTime:  startTime,
			EndTime:    endTime,
		})
	}

	// Save to database
	err := storage.AddAvailability(req.ProviderId, availabilities)
	if err != nil {
		return nil, err
	}

	return &pb.SetAvailabilityResponse{Message: "Availability set successfully"}, nil
}

func (s *ReservationService) GetAvailableSlots(ctx context.Context, req *pb.GetAvailableSlotsRequest) (*pb.GetAvailableSlotsResponse, error) {
	// Parse date
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	// Fetch slots from database
	slots, err := storage.GetAvailableSlots(req.ProviderId, date)
	if err != nil {
		return nil, err
	}

	// Convert database models to protobuf response
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
	// Calculate expiration time
	expiration := time.Now().Add(30 * time.Minute)

	// Reserve the slot in the database
	err := storage.ReserveSlot(req.SlotId, req.ClientId, expiration)
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
