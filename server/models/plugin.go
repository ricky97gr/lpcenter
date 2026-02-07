package models

import "time"

type Plugin struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID          string    `json:"uuid" gorm:"column:uuid;type:varchar(36);uniqueIndex"`
	ProductID     uint      `json:"productId" gorm:"column:product_id;not null"`
	Product       *Product  `json:"product" gorm:"foreignKey:ProductID"`
	VersionType   string    `json:"versionType" gorm:"column:version_type;type:varchar(20);not null"`
	Code          string    `json:"code" gorm:"column:code;type:varchar(50);not null;uniqueIndex"`
	Version       string    `json:"version" gorm:"column:version;type:varchar(20);not null"`
	Description   string    `json:"description" gorm:"column:description;type:text"`
	Author        string    `json:"author" gorm:"column:author;type:varchar(50)"`
	FilePath      string    `json:"filePath" gorm:"column:file_path;type:varchar(255);not null"`
	DownloadURL   string    `json:"downloadUrl" gorm:"column:download_url;type:varchar(255)"`
	MD5           string    `json:"md5" gorm:"column:md5;type:varchar(32)"`
	Status        string    `json:"status" gorm:"column:status;type:varchar(20);default:pending"`
	DownloadCount int       `json:"downloadCount" gorm:"-"`
	UploadedAt    time.Time `json:"uploadedAt" gorm:"column:uploaded_at;autoCreateTime"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (Plugin) TableName() string {
	return "plugins"
}

const (
	VersionTypeNormal = "normal"
	VersionTypePro    = "pro"
	VersionTypePlus   = "plus"
	VersionTypeMax    = "max"
)

const (
	PluginStatusPending   = "pending"
	PluginStatusPublished = "published"
	PluginStatusDisabled  = "disabled"
)
