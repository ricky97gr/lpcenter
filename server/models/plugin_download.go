package models

import "time"

type PluginDownload struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PluginID  uint      `json:"pluginId" gorm:"column:plugin_id;not null;index"`
	Plugin    *Plugin   `json:"plugin" gorm:"foreignKey:PluginID"`
	License   string    `json:"license" gorm:"column:license;type:varchar(255);not null"`
	IP        string    `json:"ip" gorm:"column:ip;type:varchar(45);not null"`
	UserAgent string    `json:"userAgent" gorm:"column:user_agent;type:text"`
	DownloadedAt time.Time `json:"downloadedAt" gorm:"column:downloaded_at;autoCreateTime"`
}

func (PluginDownload) TableName() string {
	return "plugin_downloads"
}
