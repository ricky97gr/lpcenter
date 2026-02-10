package models

import "time"

type DownloadTask struct {
	ID              uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	UUID            string     `json:"uuid" gorm:"column:uuid;type:varchar(36);uniqueIndex;not null"`
	PluginID        uint       `json:"pluginId" gorm:"column:plugin_id;not null;index"`
	Plugin          *Plugin    `json:"plugin" gorm:"foreignKey:PluginID"`
	License         string     `json:"license" gorm:"column:license;type:varchar(255);not null"`
	FilePath        string     `json:"filePath" gorm:"column:file_path;type:varchar(512);not null"`
	FileSize        int64      `json:"fileSize" gorm:"column:file_size;not null"`
	DownloadedBytes int64      `json:"downloadedBytes" gorm:"column:downloaded_bytes;default:0"`
	Progress        float64    `json:"progress" gorm:"column:progress;default:0"`
	Status          string     `json:"status" gorm:"column:status;type:varchar(20);not null;default:'pending'"`
	IP              string     `json:"ip" gorm:"column:ip;type:varchar(45);not null"`
	UserAgent       string     `json:"userAgent" gorm:"column:user_agent;type:text"`
	StartedAt       time.Time  `json:"startedAt" gorm:"column:started_at;autoCreateTime"`
	CompletedAt     *time.Time `json:"completedAt" gorm:"column:completed_at"`
}

func (DownloadTask) TableName() string {
	return "download_tasks"
}

const (
	DownloadTaskStatusPending     = "pending"
	DownloadTaskStatusDownloading = "downloading"
	DownloadTaskStatusCompleted   = "completed"
	DownloadTaskStatusFailed      = "failed"
)
