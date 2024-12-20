package storage

import (
	"log"

	"github.com/manueldelreal/health-reservation-system/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBHandler interface {
	First(dest interface{}, args ...interface{}) error
	Create(value interface{}) error
	Transaction(txFunc func(tx *gorm.DB) error) error
}

type GormDBHandler struct {
	DB *gorm.DB
}

func (g *GormDBHandler) First(dest interface{}, args ...interface{}) error {
	result := g.DB.First(dest, args...)
	return result.Error
}

func (g *GormDBHandler) Create(value interface{}) error {
	result := g.DB.Create(value)
	return result.Error
}

func (g *GormDBHandler) Transaction(txFunc func(tx *gorm.DB) error) error {
	return g.DB.Transaction(txFunc)
}

func (g *GormDBHandler) GetDB() *gorm.DB {
	return g.DB
}

var DB DBHandler

func ConnectDatabase(dsn string) {
	// Open SQLite database
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	handler := &GormDBHandler{DB: db}
	DB = handler

	// Run migrations
	log.Println("Checking and applying migrations...")
	err = handler.GetDB().AutoMigrate(
		&models.Provider{},
		&models.Availability{},
		&models.Slot{},
		&models.Reservation{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrations applied successfully.")
}
