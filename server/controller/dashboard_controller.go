package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ricky97gr/lpcenter/server/database"
	"github.com/ricky97gr/lpcenter/server/models"
	"github.com/ricky97gr/lpcenter/server/response"
	"github.com/ricky97gr/lpcenter/server/utils"
)

type DashboardStats struct {
	TotalProducts    int64 `json:"totalProducts"`
	TotalLicenseTypes int64 `json:"totalLicenseTypes"`
	TotalPlugins     int64 `json:"totalPlugins"`
	TotalLicenses    int64 `json:"totalLicenses"`
	PendingLicenses  int64 `json:"pendingLicenses"`
	ApprovedLicenses int64 `json:"approvedLicenses"`
	RejectedLicenses int64 `json:"rejectedLicenses"`
	TotalUsers       int64 `json:"totalUsers"`
	ActiveUsers      int64 `json:"activeUsers"`
	TotalDownloads   int64 `json:"totalDownloads"`
}

type RecentActivity struct {
	Type      string `json:"type"`
	Message   string `json:"message"`
	CreatedAt string `json:"createdAt"`
}

func GetDashboardStats(c *gin.Context) {
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("GetDashboardStats database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var stats DashboardStats

	db.Model(&models.Product{}).Count(&stats.TotalProducts)
	db.Model(&models.LicenseType{}).Count(&stats.TotalLicenseTypes)
	db.Model(&models.Plugin{}).Count(&stats.TotalPlugins)
	db.Model(&models.License{}).Count(&stats.TotalLicenses)
	db.Model(&models.License{}).Where("status = ?", "pending").Count(&stats.PendingLicenses)
	db.Model(&models.License{}).Where("status = ?", "approved").Count(&stats.ApprovedLicenses)
	db.Model(&models.License{}).Where("status = ?", "rejected").Count(&stats.RejectedLicenses)
	db.Model(&models.User{}).Count(&stats.TotalUsers)
	db.Model(&models.User{}).Where("status = ?", "active").Count(&stats.ActiveUsers)
	db.Model(&models.PluginDownload{}).Count(&stats.TotalDownloads)

	utils.Logger.Infow("GetDashboardStats success")
	response.Success(c, stats, 1)
}

func GetRecentLicenses(c *gin.Context) {
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("GetRecentLicenses database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var licenses []models.License
	if err := db.Order("created_at DESC").Limit(10).Find(&licenses).Error; err != nil {
		utils.Logger.Errorw("GetRecentLicenses failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "获取最近授权失败")
		return
	}

	utils.Logger.Infow("GetRecentLicenses success", "count", len(licenses))
	response.Success(c, licenses, int64(len(licenses)))
}

func GetRecentPlugins(c *gin.Context) {
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("GetRecentPlugins database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var plugins []models.Plugin
	if err := db.Order("created_at DESC").Limit(10).Find(&plugins).Error; err != nil {
		utils.Logger.Errorw("GetRecentPlugins failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "获取最近插件失败")
		return
	}

	utils.Logger.Infow("GetRecentPlugins success", "count", len(plugins))
	response.Success(c, plugins, int64(len(plugins)))
}