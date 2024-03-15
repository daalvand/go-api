package db

import (
	"fmt"

	"github.com/daalvand/go-api/models"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(db *gorm.DB) *Database {
	return &Database{DB: db}
}

func InitializeDB(dialector gorm.Dialector, config *gorm.Config) (*Database, error) {
	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}
	db = db.Debug()
	return NewDatabase(db), nil
}

func (db *Database) MigrateTables() error {
	err := db.DB.AutoMigrate(&models.User{}, &models.Event{}, &models.Registration{})
	if err != nil {
		return fmt.Errorf("failed to auto-migrate database: %w", err)
	}
	return nil
}
