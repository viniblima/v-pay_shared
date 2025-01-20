package healthcheck

import (
	"log"

	"gorm.io/gorm"
)

type HealthDatabase interface {
	CheckDatabase(dsn string) error
}

type healthDatabase struct {
	db *gorm.DB
}

func (h *healthDatabase) CheckDatabase(dsn string) error {

	sqlDB, err := h.db.DB()
	if err != nil {
		log.Fatalf("error on access connection to database: %v", err)
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("error on access connect to database: %v", err)
		return err
	}

	return nil
}

func NewHealthDatabase(db *gorm.DB) HealthDatabase {
	return &healthDatabase{
		db: db,
	}
}
