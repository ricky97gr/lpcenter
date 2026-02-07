package database

import (
	"github.com/ricky97gr/lpcenter/server/models"
	"gorm.io/gorm"
)

func Start() error {
	db, err := GetDB()
	if err != nil {
		return err
	}

	if err := AutoMigrate(db); err != nil {
		return err
	}

	if err := createDefaultUser(db); err != nil {
		return err
	}

	return nil
}

func createDefaultUser(db *gorm.DB) error {
	var count int64
	db.Model(&models.User{}).Count(&count)

	if count == 0 {
		defaultUser := models.User{
			UUID:     "default-admin-uuid",
			Username: "admin",
			Email:    "admin@lpcenter.com",
			Password: "admin123",
			Status:   models.UserStatusActive,
		}

		if err := db.Create(&defaultUser).Error; err != nil {
			return err
		}
	}

	return nil
}
