package repo

import (
	"log"

	"github.com/logservice/internal/model"
	"github.com/logservice/internal/utils"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
	users := []struct {
		username string
		password string
		role     string
	}{
		{"admin", "admin123", "admin"},
		{"viewer", "viewer123", "viewer"},
	}

	for _, u := range users {
		var existing model.User

		err := db.Where("username = ?", u.username).First(&existing).Error

		if err == nil {
			// user already exists → skip
			continue
		}

		hash, err := utils.HashPassword(u.password)
		if err != nil {
			log.Println("hash error:", err)
			continue
		}

		user := model.User{
			Username:     u.username,
			PasswordHash: hash,
			Role:         u.role,
		}

		if err := db.Create(&user).Error; err != nil {
			log.Println("insert error:", err)
		}
	}
}
