package models

import "time"

type Product struct {
	UUID        string    `json:"uuid" gorm:"primaryKey;column:uuid;type:varchar(36);uniqueIndex;not null"`
	Name        string    `json:"name" gorm:"column:name;type:varchar(100);not null"`
	Description string    `json:"description" gorm:"column:description;type:text"`
	Code        string    `json:"code" gorm:"column:code;type:varchar(50);not null;uniqueIndex"`
	Status      string    `json:"status" gorm:"column:status;type:varchar(20);default:'active'"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at"`
	Plugins     []Plugin  `json:"plugins" gorm:"foreignKey:ProductUUID;references:UUID"`
}

func (Product) TableName() string {
	return "products"
}

const (
	ProductStatusActive   = "active"
	ProductStatusInactive = "inactive"
)
