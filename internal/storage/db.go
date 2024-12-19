package storage

import (
	"log"

	"github.com/manueldelreal/health-reservation-system/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase(dsn string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Running migrations...")
	err = DB.AutoMigrate(
		&models.Provider{},
		&models.Availability{},
		&models.Slot{},
		&models.Reservation{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database connected and migrations applied")
}
