package main

import (
	"log"
	"net/http"

	"github.com/manueldelreal/health-reservation-system/internal/services"

	pb "github.com/manueldelreal/health-reservation-system/api"
)

func main() {
	server := &services.ReservationService{}
	twirpHandler := pb.NewReservationServiceServer(server)

	mux := http.NewServeMux()
	mux.Handle(twirpHandler.PathPrefix(), twirpHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
