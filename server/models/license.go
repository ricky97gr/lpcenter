package models

import "time"

type License struct {
	UUID          string    `json:"uuid" gorm:"primaryKey;column:uuid;type:varchar(36);uniqueIndex;not null"`
	SerialNumber  string    `json:"serialNumber" gorm:"column:serial_number;type:varchar(100);not null;uniqueIndex"`
	ProductUUID   string    `json:"productUuid" gorm:"column:product_uuid;type:varchar(36);not null;index"`
	Product       *Product  `json:"product" gorm:"foreignKey:ProductUUID;references:UUID"`
	LicenseType   string    `json:"licenseType" gorm:"column:license_type;type:varchar(20);not null"`
	LicensePoints int       `json:"licensePoints" gorm:"column:license_points;not null"`
	ExpiryDate    time.Time `json:"expiryDate" gorm:"column:expiry_date;not null"`
	Remarks       string    `json:"remarks" gorm:"column:remarks;type:text"`
	LicenseString string    `json:"licenseString" gorm:"column:license_string;type:text"`
	Status        string    `json:"status" gorm:"column:status;type:varchar(20);default:'pending'"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (License) TableName() string {
	return "licenses"
}

const (
	LicenseStatusPending  = "pending"
	LicenseStatusApproved = "approved"
	LicenseStatusRejected = "rejected"
)
