package main

import (
	"log"
	"net/http"
	"time"

	"github.com/manueldelreal/health-reservation-system/internal/services"
	"github.com/manueldelreal/health-reservation-system/internal/storage"

	pb "github.com/manueldelreal/health-reservation-system/api"
)

func main() {
	// Connect to SQLite database
	dsn := "file:health_reservation.db?cache=shared&mode=rwc"
	storage.ConnectDatabase(dsn)

	// Start cleanup task for expired reservations
	go func() {
		for {
			err := storage.CleanupExpiredReservations()
			if err != nil {
				log.Printf("Failed to clean expired reservations: %v", err)
			}
			time.Sleep(1 * time.Minute) // Run every minute
		}
	}()

	// Initialize the Twirp server
	server := &services.ReservationService{}
	twirpHandler := pb.NewReservationServiceServer(server)

	mux := http.NewServeMux()
	mux.Handle(twirpHandler.PathPrefix(), twirpHandler)

	// Start the server
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
