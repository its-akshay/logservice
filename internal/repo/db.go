package repo

import (
	"fmt"
	"time"

	"github.com/logservice/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(dsn string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	maxAttempts := 5
	for i := 1; i <= maxAttempts; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			if err := db.AutoMigrate(&model.Log{}, &model.User{}); err != nil {
				return nil, fmt.Errorf("migration failed: %w", err)
			}
			return db, nil
		}
		fmt.Printf("DB connection failed (attempt %d/%d): %v\n", i, maxAttempts, err)
		time.Sleep(2 * time.Second)
	}
	return nil, fmt.Errorf("could not connect to DB after %d attempts: %w", maxAttempts, err)
}
