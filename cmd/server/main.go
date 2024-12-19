package main

import (
	"log"
	"net/http"

	"github.com/manueldelreal/health-reservation-system/api"
	"github.com/manueldelreal/health-reservation-system/internal/services"
)

func main() {
	server := &services.ReservationService{}
	twirpHandler := api.NewReservationServiceServer(server)

	mux := http.NewServeMux()
	mux.Handle(twirpHandler.PathPrefix(), twirpHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
