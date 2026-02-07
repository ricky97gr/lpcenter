package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ricky97gr/lpcenter/server/database"
	"github.com/ricky97gr/lpcenter/server/models"
	"github.com/ricky97gr/lpcenter/server/response"
	"github.com/ricky97gr/lpcenter/server/utils"
)

type DownloadRequest struct {
	PluginUUID string `json:"pluginUuid" binding:"required"`
	License    string `json:"license" binding:"required"`
}

func RecordDownload(c *gin.Context) {
	var req DownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Errorw("RecordDownload parameter error", "error", err)
		response.Failed(c, http.StatusBadRequest, "参数错误")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("RecordDownload database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var plugin models.Plugin
	if err := db.Where("uuid = ?", req.PluginUUID).First(&plugin).Error; err != nil {
		utils.Logger.Errorw("RecordDownload plugin not found", "pluginUuid", req.PluginUUID, "error", err)
		response.Failed(c, http.StatusNotFound, "插件不存在")
		return
	}

	if plugin.Status != models.PluginStatusPublished {
		utils.Logger.Warnw("RecordDownload plugin not published", "pluginUuid", req.PluginUUID, "status", plugin.Status)
		response.Failed(c, http.StatusBadRequest, "插件未发布")
		return
	}

	var license models.License
	if err := db.Where("serial_number = ? AND status = ?", req.License, models.LicenseStatusApproved).First(&license).Error; err != nil {
		utils.Logger.Errorw("RecordDownload license not valid", "license", req.License, "error", err)
		response.Failed(c, http.StatusBadRequest, "许可证无效")
		return
	}

	if license.ExpiryDate.Before(time.Now()) {
		utils.Logger.Warnw("RecordDownload license expired", "license", req.License, "expiryDate", license.ExpiryDate)
		response.Failed(c, http.StatusBadRequest, "许可证已过期")
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	download := models.PluginDownload{
		PluginID:    plugin.ID,
		License:     req.License,
		IP:          clientIP,
		UserAgent:   userAgent,
		DownloadedAt: time.Now(),
	}

	if err := db.Create(&download).Error; err != nil {
		utils.Logger.Errorw("RecordDownload failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "记录下载失败")
		return
	}

	utils.Logger.Infow("RecordDownload success", "pluginId", plugin.ID, "license", req.License, "ip", clientIP)
	response.Success(c, gin.H{
		"downloadUrl": plugin.DownloadURL,
		"fileName":    plugin.FilePath,
	}, 1)
}

func GetPluginDownloadStats(c *gin.Context) {
	pluginUUID := c.Param("uuid")
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("GetPluginDownloadStats database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var plugin models.Plugin
	if err := db.Where("uuid = ?", pluginUUID).First(&plugin).Error; err != nil {
		utils.Logger.Errorw("GetPluginDownloadStats plugin not found", "pluginUuid", pluginUUID, "error", err)
		response.Failed(c, http.StatusNotFound, "插件不存在")
		return
	}

	var totalDownloads int64
	if err := db.Model(&models.PluginDownload{}).Where("plugin_id = ?", plugin.ID).Count(&totalDownloads).Error; err != nil {
		utils.Logger.Errorw("GetPluginDownloadStats count failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "获取下载统计失败")
		return
	}

	var todayDownloads int64
	today := time.Now().Truncate(24 * time.Hour)
	if err := db.Model(&models.PluginDownload{}).Where("plugin_id = ? AND downloaded_at >= ?", plugin.ID, today).Count(&todayDownloads).Error; err != nil {
		utils.Logger.Errorw("GetPluginDownloadStats today count failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "获取下载统计失败")
		return
	}

	var downloads []models.PluginDownload
	if err := db.Where("plugin_id = ?", plugin.ID).Order("downloaded_at desc").Limit(100).Find(&downloads).Error; err != nil {
		utils.Logger.Errorw("GetPluginDownloadStats find failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "获取下载记录失败")
		return
	}

	utils.Logger.Infow("GetPluginDownloadStats success", "pluginUuid", pluginUUID, "total", totalDownloads)
	response.Success(c, gin.H{
		"totalDownloads": totalDownloads,
		"todayDownloads": todayDownloads,
		"downloads":      downloads,
	}, 1)
}

func GetAllPluginDownloadStats(c *gin.Context) {
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("GetAllPluginDownloadStats database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	type PluginStats struct {
		PluginID       uint   `json:"pluginId"`
		PluginUUID     string `json:"pluginUuid"`
		PluginCode     string `json:"pluginCode"`
		PluginVersion  string `json:"pluginVersion"`
		TotalDownloads int64  `json:"totalDownloads"`
	}

	var stats []PluginStats
	query := `
		SELECT 
			p.id as plugin_id,
			p.uuid as plugin_uuid,
			p.code as plugin_code,
			p.version as plugin_version,
			COUNT(pd.id) as total_downloads
		FROM plugins p
		LEFT JOIN plugin_downloads pd ON p.id = pd.plugin_id
		GROUP BY p.id
		ORDER BY total_downloads DESC
	`

	if err := db.Raw(query).Scan(&stats).Error; err != nil {
		utils.Logger.Errorw("GetAllPluginDownloadStats query failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "获取下载统计失败")
		return
	}

	utils.Logger.Infow("GetAllPluginDownloadStats success", "count", len(stats))
	response.Success(c, stats, int64(len(stats)))
}
