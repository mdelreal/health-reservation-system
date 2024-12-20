package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"

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
	// Validate that the provider exists
	var provider models.Provider
	err := storage.DB.First(&provider, "id = ?", req.ProviderId)
	if err != nil {
		return nil, errors.New("provider not found")
	}

	availabilityMap := make(map[string]models.Availability)
	var slots []models.Slot

	for _, timeSlot := range req.TimeSlots {
		// Parse start and end times
		startTime, err := time.Parse(time.RFC3339, timeSlot.StartTime)
		if err != nil {
			return nil, errors.New("invalid start time format")
		}
		endTime, err := time.Parse(time.RFC3339, timeSlot.EndTime)
		if err != nil {
			return nil, errors.New("invalid end time format")
		}

		// Check if availability already exists for this timeframe
		var existingAvailability models.Availability
		err = storage.DB.First(&existingAvailability, "provider_id = ? AND start_time = ? AND end_time = ?", req.ProviderId, startTime, endTime)
		if err == nil {
			// Skip this time slot if availability already exists
			continue
		}

		// Create a new availability entry
		availability := models.Availability{
			ID:         generateID(),
			ProviderID: req.ProviderId,
			StartTime:  startTime,
			EndTime:    endTime,
		}
		availabilityMap[fmt.Sprintf("%s-%s", startTime, endTime)] = availability

		// Split the interval into 15-minute slots
		for t := startTime; t.Before(endTime); t = t.Add(15 * time.Minute) {
			// Check if slot already exists for this timeframe
			var existingSlot models.Slot
			err = storage.DB.First(&existingSlot, "availability_id = ? AND start_time = ? AND end_time = ?", availability.ID, t, t.Add(15*time.Minute))
			if err == nil {
				// Skip this slot if it already exists
				continue
			}

			slots = append(slots, models.Slot{
				ID:             generateID(),
				AvailabilityID: availability.ID,
				StartTime:      t,
				EndTime:        t.Add(15 * time.Minute),
				Status:         "Available",
			})
		}
	}

	// Save new availabilities and slots to the database
	availabilities := make([]models.Availability, 0, len(availabilityMap))
	for _, availability := range availabilityMap {
		availabilities = append(availabilities, availability)
	}

	if len(availabilities) > 0 || len(slots) > 0 {
		err = storage.AddAvailabilityAndSlots(req.ProviderId, availabilities, slots)
		if err != nil {
			return nil, err
		}
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
	err := storage.DB.First(&slot, "id = ? AND status = ?", req.SlotId, "Available")
	if err != nil {
		return nil, errors.New("slot is not available")
	}

	// Validate 24-hour rule
	if slot.StartTime.Before(time.Now().Add(24 * time.Hour)) {
		return nil, errors.New("reservations must be made at least 24 hours in advance")
	}

	// Reserve the slot
	expiration := time.Now().Add(30 * time.Minute)
	err = storage.ReserveSlot(req.SlotId, req.ClientId, slot.ProviderID, expiration)
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

func (s *ReservationService) CreateProvider(ctx context.Context, req *pb.CreateProviderRequest) (*pb.CreateProviderResponse, error) {
	// Check if the provider already exists
	var existingProvider models.Provider
	err := storage.DB.First(&existingProvider, "id = ?", req.Id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// Return error if it is not a "record not found" error
		return nil, errors.New("failed to query provider")
	}
	if err == nil {
		// If no error, the provider already exists
		return nil, errors.New("provider already exists")
	}

	// Create a new provider
	provider := models.Provider{
		ID:   req.Id,
		Name: req.Name,
	}

	err = storage.DB.Create(&provider)
	if err != nil {
		return nil, errors.New("failed to create provider")
	}

	return &pb.CreateProviderResponse{
		Message: "Provider created successfully",
	}, nil
}

func (s *ReservationService) GetProvider(ctx context.Context, req *pb.GetProviderRequest) (*pb.GetProviderResponse, error) {
	// Retrieve the provider data
	var provider models.Provider
	err := storage.DB.First(&provider, "id = ?", req.Id)
	if err != nil {
		return nil, errors.New("provider not found")
	}

	return &pb.GetProviderResponse{
		Id:   provider.ID,
		Name: provider.Name,
	}, nil
}

func (s *ReservationService) GetReservedSlotsByProvider(ctx context.Context, req *pb.GetReservedSlotsByProviderRequest) (*pb.GetReservedSlotsByProviderResponse, error) {
	// Parse the optional date argument
	var date *time.Time
	if req.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return nil, errors.New("invalid date format")
		}
		date = &parsedDate
	}

	// Query for reservations
	reservations, err := storage.GetReservationsByProvider(req.ProviderId, date)
	if err != nil {
		return nil, err
	}

	// Convert to protobuf response
	var pbReservations []*pb.ReservationDetails
	for _, reservation := range reservations {
		pbReservations = append(pbReservations, &pb.ReservationDetails{
			ReservationId: reservation.ID,
			ClientId:      reservation.ClientID,
			ProviderId:    req.ProviderId,
			Status:        reservation.Status,
			StartTime:     reservation.StartTime.Format(time.RFC3339),
			EndTime:       reservation.EndTime.Format(time.RFC3339),
		})
	}

	return &pb.GetReservedSlotsByProviderResponse{Reservations: pbReservations}, nil
}

func (s *ReservationService) GetReservedSlotsByClient(ctx context.Context, req *pb.GetReservedSlotsByClientRequest) (*pb.GetReservedSlotsByClientResponse, error) {
	// Parse the optional date argument
	var date *time.Time
	if req.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return nil, errors.New("invalid date format")
		}
		date = &parsedDate
	}

	// Query for reservations
	reservations, err := storage.GetReservationsByClient(req.ClientId, date)
	if err != nil {
		return nil, err
	}

	// Convert to protobuf response
	var pbReservations []*pb.ReservationDetails
	for _, reservation := range reservations {
		pbReservations = append(pbReservations, &pb.ReservationDetails{
			ReservationId: reservation.ID,
			ClientId:      req.ClientId,
			ProviderId:    reservation.ProviderID,
			Status:        reservation.Status,
			StartTime:     reservation.StartTime.Format(time.RFC3339),
			EndTime:       reservation.EndTime.Format(time.RFC3339),
		})
	}

	return &pb.GetReservedSlotsByClientResponse{Reservations: pbReservations}, nil
}
