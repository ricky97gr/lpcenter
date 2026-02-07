package models

import "time"

type License struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID           string    `json:"uuid" gorm:"column:uuid;type:varchar(36);uniqueIndex"`
	SerialNumber   string    `json:"serialNumber" gorm:"column:serial_number;type:varchar(100);not null;uniqueIndex"`
	ProductID      uint      `json:"productId" gorm:"column:product_id;not null"`
	Product        *Product  `json:"product" gorm:"foreignKey:ProductID"`
	Version        string    `json:"version" gorm:"column:version;type:varchar(20);not null"`
	ExpiryDate     time.Time `json:"expiryDate" gorm:"column:expiry_date;not null"`
	Remarks        string    `json:"remarks" gorm:"column:remarks;type:text"`
	LicenseString  string    `json:"licenseString" gorm:"column:license_string;type:text"`
	Status         string    `json:"status" gorm:"column:status;type:varchar(20);default:'pending'"`
	CreatedAt      time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt      time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (License) TableName() string {
	return "licenses"
}

const (
	LicenseStatusPending  = "pending"
	LicenseStatusApproved = "approved"
	LicenseStatusRejected = "rejected"
)
