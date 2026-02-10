package models

import "time"

type Plugin struct {
	ID                uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID              string    `json:"uuid" gorm:"column:uuid;type:varchar(36);uniqueIndex"`
	Name              string    `json:"name" gorm:"column:name;type:varchar(50);not null"`
	ProductID         uint      `json:"productId" gorm:"column:product_id;not null"`
	Product           *Product  `json:"product" gorm:"foreignKey:ProductID"`
	LicenseType       string    `json:"licenseType" gorm:"column:license_type;type:varchar(20);not null"`
	Code              string    `json:"code" gorm:"column:code;type:varchar(50);not null;uniqueIndex"`
	Version           string    `json:"version" gorm:"column:version;type:varchar(20);not null"`
	Description       string    `json:"description" gorm:"column:description;type:text"`
	Tips              string    `json:"tips" gorm:"column:tips;type:text"`
	Author            string    `json:"author" gorm:"column:author;type:varchar(50)"`
	Contact           string    `json:"contact" gorm:"column:contact;type:varchar(100)"`
	FilePath          string    `json:"filePath" gorm:"column:file_path;type:varchar(255);not null"`
	DownloadURL       string    `json:"downloadUrl" gorm:"column:download_url;type:varchar(255)"`
	MD5               string    `json:"md5" gorm:"column:md5;type:varchar(32)"`
	Status            string    `json:"status" gorm:"column:status;type:varchar(20);default:pending"`
	Signed            bool      `json:"signed" gorm:"column:signed;type:boolean;default:false"`
	DownloadCount     int       `json:"downloadCount" gorm:"-"`
	UploadedAt        time.Time `json:"uploadedAt" gorm:"column:uploaded_at;autoCreateTime"`
	MiniServerVersion string    `json:"miniServerVersion" gorm:"column:mini_server_version;type:varchar(20)"`
	MiniClientVersion string    `json:"miniClientVersion" gorm:"column:mini_client_version;type:varchar(20)"`
	CreatedAt         time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt         time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (Plugin) TableName() string {
	return "plugins"
}

const (
	LicenseTypeNormal = "normal"
	LicenseTypePro    = "pro"
	LicenseTypePlus   = "plus"
	LicenseTypeMax    = "max"
)

const (
	PluginStatusPending   = "pending"
	PluginStatusPublished = "published"
	PluginStatusSigned    = "signed"
	PluginStatusDisabled  = "disabled"
)
