package storage

import (
	"log"
	"os"

	"github.com/manueldelreal/health-reservation-system/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase(dsn string) {
	// Check if the database file exists
	_, err := os.Stat(dsn)
	isNewDB := os.IsNotExist(err)

	// Open the database
	DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if isNewDB {
		log.Println("Database file does not exist. Creating new database and running migrations...")
		err = DB.AutoMigrate(
			&models.Provider{},
			&models.Availability{},
			&models.Slot{},
			&models.Reservation{},
		)
		if err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}
		log.Println("Database created and migrations applied.")
	} else {
		log.Println("Database already exists. Skipping migrations.")
	}
}
