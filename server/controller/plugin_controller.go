package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ricky97gr/lpcenter/server/database"
	"github.com/ricky97gr/lpcenter/server/models"
	"github.com/ricky97gr/lpcenter/server/pagination"
	"github.com/ricky97gr/lpcenter/server/response"
	"github.com/ricky97gr/lpcenter/server/utils"
)

type PluginRequest struct {
	ProductID         uint   `json:"productId" binding:"required"`
	VersionType       string `json:"versionType" binding:"required"`
	Code              string `json:"code" binding:"required"`
	Name              string `json:"name" binding:"required"`
	Version           string `json:"version" binding:"required"`
	MiniServerVersion string `json:"miniServerVersion"`
	MiniClientVersion string `json:"miniClientVersion"`
	Description       string `json:"description"`
	Tips              string `json:"tips"`
	Author            string `json:"author"`
	FilePath          string `json:"filePath" binding:"required"`
	DownloadURL       string `json:"downloadUrl"`
	Status            string `json:"status"`
}

func CreatePlugin(c *gin.Context) {
	var req PluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Errorw("CreatePlugin parameter error", "error", err)
		response.Failed(c, http.StatusBadRequest, "参数错误")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("CreatePlugin database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var product models.Product
	if err := db.First(&product, req.ProductID).Error; err != nil {
		utils.Logger.Errorw("CreatePlugin product not found", "productId", req.ProductID, "error", err)
		response.Failed(c, http.StatusNotFound, "产品不存在")
		return
	}

	plugin := models.Plugin{
		UUID:              uuid.New().String(),
		ProductID:         req.ProductID,
		VersionType:       req.VersionType,
		Code:              req.Code,
		Name:              req.Name,
		Version:           req.Version,
		MiniServerVersion: req.MiniServerVersion,
		MiniClientVersion: req.MiniClientVersion,
		Description:       req.Description,
		Tips:              req.Tips,
		Author:            req.Author,
		FilePath:          req.FilePath,
		DownloadURL:       req.DownloadURL,
		Status:            models.PluginStatusPending,
	}

	if err := db.Create(&plugin).Error; err != nil {
		utils.Logger.Errorw("CreatePlugin failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "创建插件失败")
		return
	}

	utils.Logger.Infow("CreatePlugin success", "pluginId", plugin.ID, "name", plugin.Code)
	response.Success(c, plugin, 1)
}

func GetAllPlugins(c *gin.Context) {
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("GetAllPlugins database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	pageQuery, err := pagination.GetPageQuery(c)
	if err != nil {
		utils.Logger.Errorw("GetAllPlugins page query error", "error", err)
		response.Failed(c, http.StatusBadRequest, "分页参数错误")
		return
	}

	var plugins []models.Plugin
	var total int64

	countDB := db.Model(&models.Plugin{})
	for _, cond := range pageQuery.Conditions {
		countDB = countDB.Scopes(pagination.QueryFilter(cond.Field, cond.Value, cond.Operation))
	}
	if err := countDB.Count(&total).Error; err != nil {
		utils.Logger.Errorw("GetAllPlugins count failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "获取插件列表失败")
		return
	}

	queryDB := db.Preload("Product").Scopes(pagination.ParseQuery(pageQuery))
	if err := queryDB.Find(&plugins).Error; err != nil {
		utils.Logger.Errorw("GetAllPlugins find failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "获取插件列表失败")
		return
	}

	for i := range plugins {
		var downloadCount int64
		db.Model(&models.PluginDownload{}).Where("plugin_id = ?", plugins[i].ID).Count(&downloadCount)
		plugins[i].DownloadCount = int(downloadCount)
	}

	utils.Logger.Infow("GetAllPlugins success", "count", total, "page", pageQuery.Page, "pageSize", pageQuery.PageSize)
	response.Success(c, plugins, total)
}

func GetPublicPlugins(c *gin.Context) {
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("GetPublicPlugins database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	pageQuery, err := pagination.GetPageQuery(c)
	if err != nil {
		utils.Logger.Errorw("GetPublicPlugins page query error", "error", err)
		response.Failed(c, http.StatusBadRequest, "分页参数错误")
		return
	}

	var plugins []models.Plugin
	var total int64

	countDB := db.Model(&models.Plugin{}).Where("status = ?", models.PluginStatusPublished)
	for _, cond := range pageQuery.Conditions {
		countDB = countDB.Scopes(pagination.QueryFilter(cond.Field, cond.Value, cond.Operation))
	}
	if err := countDB.Count(&total).Error; err != nil {
		utils.Logger.Errorw("GetPublicPlugins count failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "获取插件列表失败")
		return
	}

	queryDB := db.Preload("Product").Where("status = ?", models.PluginStatusPublished).Scopes(pagination.ParseQuery(pageQuery))
	if err := queryDB.Find(&plugins).Error; err != nil {
		utils.Logger.Errorw("GetPublicPlugins find failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "获取插件列表失败")
		return
	}

	for i := range plugins {
		var downloadCount int64
		db.Model(&models.PluginDownload{}).Where("plugin_id = ?", plugins[i].ID).Count(&downloadCount)
		plugins[i].DownloadCount = int(downloadCount)
	}

	utils.Logger.Infow("GetPublicPlugins success", "count", total, "page", pageQuery.Page, "pageSize", pageQuery.PageSize)
	response.Success(c, plugins, total)
}

func GetPlugin(c *gin.Context) {
	id := c.Param("id")
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("GetPlugin database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var plugin models.Plugin

	if err := db.Preload("Product").First(&plugin, id).Error; err != nil {
		utils.Logger.Errorw("GetPlugin not found", "id", id, "error", err)
		response.Failed(c, http.StatusNotFound, "插件不存在")
		return
	}

	utils.Logger.Infow("GetPlugin success", "pluginId", plugin.ID)
	response.Success(c, plugin, 1)
}

func UpdatePlugin(c *gin.Context) {
	id := c.Param("id")
	var req PluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Errorw("UpdatePlugin parameter error", "id", id, "error", err)
		response.Failed(c, http.StatusBadRequest, "参数错误")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("UpdatePlugin database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var plugin models.Plugin

	if err := db.First(&plugin, id).Error; err != nil {
		utils.Logger.Errorw("UpdatePlugin not found", "id", id, "error", err)
		response.Failed(c, http.StatusNotFound, "插件不存在")
		return
	}

	var product models.Product
	if err := db.First(&product, req.ProductID).Error; err != nil {
		utils.Logger.Errorw("UpdatePlugin product not found", "productId", req.ProductID, "error", err)
		response.Failed(c, http.StatusNotFound, "产品不存在")
		return
	}

	plugin.ProductID = req.ProductID
	plugin.VersionType = req.VersionType
	plugin.Code = req.Code
	plugin.Name = req.Name
	plugin.Version = req.Version
	plugin.MiniServerVersion = req.MiniServerVersion
	plugin.MiniClientVersion = req.MiniClientVersion
	plugin.Description = req.Description
	plugin.Tips = req.Tips
	plugin.Author = req.Author
	plugin.FilePath = req.FilePath
	plugin.DownloadURL = req.DownloadURL
	if req.Status != "" {
		plugin.Status = req.Status
	}

	if err := db.Save(&plugin).Error; err != nil {
		utils.Logger.Errorw("UpdatePlugin failed", "id", id, "error", err)
		response.Failed(c, http.StatusInternalServerError, "更新插件失败")
		return
	}

	utils.Logger.Infow("UpdatePlugin success", "pluginId", plugin.ID)
	response.Success(c, plugin, 1)
}

func DeletePlugin(c *gin.Context) {
	id := c.Param("id")
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("DeletePlugin database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	if err := db.Delete(&models.Plugin{}, id).Error; err != nil {
		utils.Logger.Errorw("DeletePlugin failed", "id", id, "error", err)
		response.Failed(c, http.StatusInternalServerError, "删除插件失败")
		return
	}

	utils.Logger.Infow("DeletePlugin success", "pluginId", id)
	response.Success(c, nil, 0)
}

func UpdatePluginStatus(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Errorw("UpdatePluginStatus parameter error", "id", id, "error", err)
		response.Failed(c, http.StatusBadRequest, "参数错误")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("UpdatePluginStatus database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var plugin models.Plugin

	if err := db.First(&plugin, id).Error; err != nil {
		utils.Logger.Errorw("UpdatePluginStatus not found", "id", id, "error", err)
		response.Failed(c, http.StatusNotFound, "插件不存在")
		return
	}

	plugin.Status = req.Status

	if err := db.Save(&plugin).Error; err != nil {
		utils.Logger.Errorw("UpdatePluginStatus failed", "id", id, "error", err)
		response.Failed(c, http.StatusInternalServerError, "更新插件状态失败")
		return
	}

	utils.Logger.Infow("UpdatePluginStatus success", "pluginId", plugin.ID, "status", plugin.Status)
	response.Success(c, plugin, 1)
}
