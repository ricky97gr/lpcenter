package models

import "time"

type Product struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID        string    `json:"uuid" gorm:"column:uuid;type:varchar(36);uniqueIndex"`
	Name        string    `json:"name" gorm:"column:name;type:varchar(100);not null"`
	Description string    `json:"description" gorm:"column:description;type:text"`
	Code        string    `json:"code" gorm:"column:code;type:varchar(50);not null;uniqueIndex"`
	Status      string    `json:"status" gorm:"column:status;type:varchar(20);default:'active'"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at"`
	Plugins     []Plugin  `json:"plugins" gorm:"foreignKey:ProductID"`
}

func (Product) TableName() string {
	return "products"
}

const (
	ProductStatusActive   = "active"
	ProductStatusInactive = "inactive"
)
