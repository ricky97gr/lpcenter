package database

import (
	"github.com/ricky97gr/lpcenter/server/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Version{},
		&models.Product{},
		&models.Plugin{},
		&models.License{},
		&models.PluginDownload{},
	)
}
