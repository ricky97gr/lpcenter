package models

import "time"

type PluginDownload struct {
	UUID         string    `json:"uuid" gorm:"primaryKey;column:uuid;type:varchar(36);uniqueIndex;not null"`
	PluginUUID   string    `json:"pluginUuid" gorm:"column:plugin_uuid;type:varchar(36);not null;index"`
	Plugin       *Plugin   `json:"plugin" gorm:"foreignKey:PluginUUID;references:UUID"`
	License      string    `json:"license" gorm:"column:license;type:varchar(255);not null"`
	IP           string    `json:"ip" gorm:"column:ip;type:varchar(45);not null"`
	UserAgent    string    `json:"userAgent" gorm:"column:user_agent;type:text"`
	DownloadedAt time.Time `json:"downloadedAt" gorm:"column:downloaded_at;autoCreateTime"`
}

func (PluginDownload) TableName() string {
	return "plugin_downloads"
}
