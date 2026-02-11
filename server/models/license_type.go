package models

import "time"

type LicenseType struct {
	UUID      string    `json:"uuid" gorm:"column:uuid;type:varchar(36);uniqueIndex"`
	Name      string    `json:"name" gorm:"column:name;type:varchar(50);not null;uniqueIndex"`
	Code      string    `json:"code" gorm:"column:code;type:varchar(50);not null;uniqueIndex"`
	IsPaid    bool      `json:"isPaid" gorm:"column:is_paid;type:boolean;default:false"`
	CreatedBy string    `json:"createdBy" gorm:"column:created_by;type:varchar(100)"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (LicenseType) TableName() string {
	return "license_types"
}
