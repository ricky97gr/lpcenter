package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ricky97gr/lpcenter/server/database"
	"github.com/ricky97gr/lpcenter/server/models"
	"github.com/ricky97gr/lpcenter/server/pagination"
	"github.com/ricky97gr/lpcenter/server/response"
	"github.com/ricky97gr/lpcenter/server/utils"
)

type PluginRequest struct {
	ProductUUID       string `json:"productUuid" binding:"required"`
	LicenseType       string `json:"licenseType" binding:"required"`
	Name              string `json:"name" binding:"required"`
	Version           string `json:"version" binding:"required"`
	MiniServerVersion string `json:"miniServerVersion"`
	MiniClientVersion string `json:"miniClientVersion"`
	Description       string `json:"description"`
	Tips              string `json:"tips"`
	Author            string `json:"author"`
	Contact           string `json:"contact"`
	FileName          string `json:"fileName" binding:"required"`
	FilePath          string `json:"filePath"`
	DownloadURL       string `json:"downloadUrl"`
	Status            string `json:"status"`
}

type DownloadRequest struct {
	UUID    string `json:"uuid" binding:"required"`
	License string `json:"license" binding:"required"`
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
	if err := db.First(&product, "uuid = ?", req.ProductUUID).Error; err != nil {
		utils.Logger.Errorw("CreatePlugin product not found", "productUuid", req.ProductUUID, "error", err)
		response.Failed(c, http.StatusNotFound, "产品不存在")
		return
	}

	// 生成文件存储路径
	pluginDir := "./uploads/plugins/"
	filePath := pluginDir + req.FileName

	// 确保目录存在
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		utils.Logger.Errorw("CreatePlugin create directory failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "创建插件目录失败")
		return
	}

	// 生成插件编号
	productCode := "plugin"
	if product.Code != "" {
		productCode = product.Code
	}
	code := fmt.Sprintf("%s_%s_%d", productCode, req.LicenseType, time.Now().Unix())

	// 生成下载链接
	downloadURL := "http://localhost:9092/"

	plugin := models.Plugin{
		UUID:              uuid.New().String(),
		ProductUUID:       req.ProductUUID,
		LicenseType:       req.LicenseType,
		Code:              code,
		Name:              req.Name,
		Version:           req.Version,
		MiniServerVersion: req.MiniServerVersion,
		MiniClientVersion: req.MiniClientVersion,
		Description:       req.Description,
		Tips:              req.Tips,
		Author:            req.Author,
		Contact:           req.Contact,
		FilePath:          filePath,
		DownloadURL:       downloadURL,
		Status:            models.PluginStatusPending,
	}

	if err := db.Create(&plugin).Error; err != nil {
		utils.Logger.Errorw("CreatePlugin failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "创建插件失败")
		return
	}

	utils.Logger.Infow("CreatePlugin success", "pluginId", plugin.UUID, "name", plugin.Code)
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
		db.Model(&models.PluginDownload{}).Where("plugin_uuid = ?", plugins[i].UUID).Count(&downloadCount)
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
		db.Model(&models.PluginDownload{}).Where("plugin_uuid = ?", plugins[i].UUID).Count(&downloadCount)
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

	if err := db.Preload("Product").First(&plugin, "uuid = ?", id).Error; err != nil {
		utils.Logger.Errorw("GetPlugin not found", "id", id, "error", err)
		response.Failed(c, http.StatusNotFound, "插件不存在")
		return
	}

	utils.Logger.Infow("GetPlugin success", "pluginId", plugin.UUID)
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

	if err := db.First(&plugin, "uuid = ?", id).Error; err != nil {
		utils.Logger.Errorw("UpdatePlugin not found", "id", id, "error", err)
		response.Failed(c, http.StatusNotFound, "插件不存在")
		return
	}

	var product models.Product
	if err := db.First(&product, "uuid = ?", req.ProductUUID).Error; err != nil {
		utils.Logger.Errorw("UpdatePlugin product not found", "productUuid", req.ProductUUID, "error", err)
		response.Failed(c, http.StatusNotFound, "产品不存在")
		return
	}

	plugin.ProductUUID = req.ProductUUID
	plugin.LicenseType = req.LicenseType
	// 保持原来的插件编号
	plugin.Name = req.Name
	plugin.Version = req.Version
	plugin.MiniServerVersion = req.MiniServerVersion
	plugin.MiniClientVersion = req.MiniClientVersion
	plugin.Description = req.Description
	plugin.Tips = req.Tips
	plugin.Author = req.Author
	plugin.Contact = req.Contact
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

	utils.Logger.Infow("UpdatePlugin success", "pluginId", plugin.UUID)
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

	// 先获取插件信息，用于删除文件
	var plugin models.Plugin
	if err := db.First(&plugin, "uuid = ?", id).Error; err == nil {
		// 删除对应的文件
		if plugin.FilePath != "" {
			if err := os.Remove(plugin.FilePath); err != nil {
				// 记录错误但不阻止删除操作
				utils.Logger.Errorw("DeletePlugin remove file failed", "filePath", plugin.FilePath, "error", err)
			}
		}
	}

	if err := db.Where("uuid = ?", id).Delete(&models.Plugin{}).Error; err != nil {
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

	if err := db.First(&plugin, "uuid = ?", id).Error; err != nil {
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

	utils.Logger.Infow("UpdatePluginStatus success", "pluginId", plugin.UUID, "status", plugin.Status)
	response.Success(c, plugin, 1)
}

func SignPlugin(c *gin.Context) {
	id := c.Param("id")
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("SignPlugin database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var plugin models.Plugin

	if err := db.Preload("Product").First(&plugin, "uuid = ?", id).Error; err != nil {
		utils.Logger.Errorw("SignPlugin not found", "id", id, "error", err)
		response.Failed(c, http.StatusNotFound, "插件不存在")
		return
	}

	// 确保插件编号存在
	if plugin.Code == "" {
		productCode := "plugin"
		if plugin.Product != nil && plugin.Product.Code != "" {
			productCode = plugin.Product.Code
		}
		plugin.Code = fmt.Sprintf("%s_%s_%d", productCode, plugin.LicenseType, time.Now().Unix())
	}

	// 确保文件存储目录存在
	pluginDir := "./uploads/plugins/"
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		utils.Logger.Errorw("SignPlugin create directory failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "创建插件目录失败")
		return
	}

	// 处理文件路径
	if plugin.FilePath != "" && !strings.HasPrefix(plugin.FilePath, pluginDir) {
		// 如果文件路径不是绝对路径，转换为绝对路径
		if !strings.HasPrefix(plugin.FilePath, "/") {
			plugin.FilePath = pluginDir + plugin.FilePath
		}
	}

	plugin.Signed = true
	plugin.Status = models.PluginStatusSigned

	if err := db.Save(&plugin).Error; err != nil {
		utils.Logger.Errorw("SignPlugin failed", "id", id, "error", err)
		response.Failed(c, http.StatusInternalServerError, "签名插件失败")
		return
	}

	utils.Logger.Infow("SignPlugin success", "pluginId", plugin.UUID, "code", plugin.Code, "filePath", plugin.FilePath)
	response.Success(c, plugin, 1)
}

func PublishPlugin(c *gin.Context) {
	id := c.Param("id")
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("PublishPlugin database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var plugin models.Plugin

	if err := db.First(&plugin, "uuid = ?", id).Error; err != nil {
		utils.Logger.Errorw("PublishPlugin not found", "id", id, "error", err)
		response.Failed(c, http.StatusNotFound, "插件不存在")
		return
	}

	plugin.Status = models.PluginStatusPublished

	if err := db.Save(&plugin).Error; err != nil {
		utils.Logger.Errorw("PublishPlugin failed", "id", id, "error", err)
		response.Failed(c, http.StatusInternalServerError, "发布插件失败")
		return
	}

	utils.Logger.Infow("PublishPlugin success", "pluginId", plugin.UUID)
	response.Success(c, plugin, 1)
}

func DisablePlugin(c *gin.Context) {
	id := c.Param("id")
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("DisablePlugin database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var plugin models.Plugin

	if err := db.First(&plugin, "uuid = ?", id).Error; err != nil {
		utils.Logger.Errorw("DisablePlugin not found", "id", id, "error", err)
		response.Failed(c, http.StatusNotFound, "插件不存在")
		return
	}

	plugin.Status = models.PluginStatusDisabled

	if err := db.Save(&plugin).Error; err != nil {
		utils.Logger.Errorw("DisablePlugin failed", "id", id, "error", err)
		response.Failed(c, http.StatusInternalServerError, "停用插件失败")
		return
	}

	utils.Logger.Infow("DisablePlugin success", "pluginId", plugin.UUID)
	response.Success(c, plugin, 1)
}

func UploadPlugin(c *gin.Context) {
	// 解析表单数据
	productUUID := c.PostForm("productUuid")
	if productUUID == "" {
		response.Failed(c, http.StatusBadRequest, "产品UUID参数错误")
		return
	}

	licenseType := c.PostForm("licenseType")
	name := c.PostForm("name")
	version := c.PostForm("version")
	miniServerVersion := c.PostForm("miniServerVersion")
	miniClientVersion := c.PostForm("miniClientVersion")
	description := c.PostForm("description")
	tips := c.PostForm("tips")
	author := c.PostForm("author")
	contact := c.PostForm("contact")

	// 检查必填参数
	if licenseType == "" || name == "" || version == "" {
		response.Failed(c, http.StatusBadRequest, "缺少必填参数")
		return
	}

	// 处理文件上传
	file, err := c.FormFile("file")
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "文件上传失败")
		return
	}

	// 确保目录存在
	pluginDir := "./uploads/plugins/"
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		response.Failed(c, http.StatusInternalServerError, "创建插件目录失败")
		return
	}

	// 保存文件
	filename := file.Filename
	filePath := pluginDir + filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		response.Failed(c, http.StatusInternalServerError, "保存文件失败")
		return
	}

	// 连接数据库
	db, err := database.GetDB()
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	// 检查产品是否存在
	var product models.Product
	if err := db.First(&product, "uuid = ?", productUUID).Error; err != nil {
		response.Failed(c, http.StatusNotFound, "产品不存在")
		return
	}

	// 生成插件编号
	productCode := "plugin"
	if product.Code != "" {
		productCode = product.Code
	}
	code := fmt.Sprintf("%s_%s_%d", productCode, licenseType, time.Now().Unix())

	// 生成下载链接
	downloadURL := fmt.Sprintf("http://localhost:9092/download/%s", filename)

	// 创建插件记录
	plugin := models.Plugin{
		UUID:              uuid.New().String(),
		ProductUUID:       productUUID,
		LicenseType:       licenseType,
		Code:              code,
		Name:              name,
		Version:           version,
		MiniServerVersion: miniServerVersion,
		MiniClientVersion: miniClientVersion,
		Description:       description,
		Tips:              tips,
		Author:            author,
		Contact:           contact,
		FilePath:          filePath,
		DownloadURL:       downloadURL,
		Status:            models.PluginStatusPending,
	}

	if err := db.Create(&plugin).Error; err != nil {
		response.Failed(c, http.StatusInternalServerError, "创建插件失败")
		return
	}

	response.Success(c, plugin, 1)
}

// 创建下载任务接口
func CreateDownloadTask(c *gin.Context) {
	var req DownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Errorw("CreateDownloadTask parameter error", "error", err)
		response.Failed(c, http.StatusBadRequest, "参数错误")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("CreateDownloadTask database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	// 查找插件
	var plugin models.Plugin
	fmt.Println("find plugin, uuid: ", req.UUID)
	if err := db.First(&plugin, req.UUID).Error; err != nil {
		response.Failed(c, http.StatusNotFound, "插件不存在")
		return
	}

	// 验证插件状态
	if plugin.Status != models.PluginStatusPublished {
		response.Failed(c, http.StatusBadRequest, "插件未发布")
		return
	}

	// 验证许可证
	var license models.License
	if err := db.Where("license_string = ? AND status = ?", req.License, models.LicenseStatusApproved).First(&license).Error; err != nil {
		response.Failed(c, http.StatusBadRequest, "无效的许可证")
		return
	}

	// 验证许可证是否过期
	if license.ExpiryDate.Before(time.Now()) {
		response.Failed(c, http.StatusBadRequest, "许可证已过期")
		return
	}

	// 验证许可证类型是否与插件类型匹配
	if license.LicenseType != plugin.LicenseType {
		response.Failed(c, http.StatusBadRequest, "许可证类型与插件类型不匹配")
		return
	}

	// 验证许可证产品是否与插件产品匹配
	if license.ProductUUID != plugin.ProductUUID {
		response.Failed(c, http.StatusBadRequest, "许可证产品与插件产品不匹配")
		return
	}

	// 检查是否已存在相同的完成任务
	var existingTask models.DownloadTask
	if err := db.Where("uuid = ? AND license = ? AND status = ?", req.UUID, req.License, models.DownloadTaskStatusCompleted).First(&existingTask).Error; err == nil {
		response.Failed(c, http.StatusBadRequest, "该插件已下载完成，不允许再次下载")
		return
	}

	// 检查文件是否存在
	if !strings.HasPrefix(plugin.FilePath, "./uploads/plugins/") {
		plugin.FilePath = "./uploads/plugins/" + filepath.Base(plugin.FilePath)
	}

	fileInfo, err := os.Stat(plugin.FilePath)
	if os.IsNotExist(err) {
		response.Failed(c, http.StatusNotFound, "插件文件不存在")
		return
	}

	// 创建下载任务
	taskUUID := uuid.New().String()
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	task := models.DownloadTask{
		UUID:       taskUUID,
		PluginUUID: req.UUID,
		License:    req.License,
		FilePath:   plugin.FilePath,
		FileSize:   fileInfo.Size(),
		Status:     models.DownloadTaskStatusPending,
		IP:         clientIP,
		UserAgent:  userAgent,
		StartedAt:  time.Now(),
	}

	if err := db.Create(&task).Error; err != nil {
		utils.Logger.Errorw("CreateDownloadTask failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "创建下载任务失败")
		return
	}

	// 返回任务信息
	response.Success(c, gin.H{
		"taskId":      taskUUID,
		"fileName":    filepath.Base(plugin.FilePath),
		"fileSize":    fileInfo.Size(),
		"downloadUrl": fmt.Sprintf("http://localhost:8082/download/file?taskId=%s", taskUUID),
	}, 1)
}

// 查询下载状态接口
func GetDownloadStatus(c *gin.Context) {
	taskUUID := c.Param("taskId")

	if taskUUID == "" {
		response.Failed(c, http.StatusBadRequest, "任务ID不能为空")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var task models.DownloadTask
	if err := db.Preload("Plugin").Where("uuid = ?", taskUUID).First(&task).Error; err != nil {
		response.Failed(c, http.StatusNotFound, "下载任务不存在")
		return
	}

	response.Success(c, task, 1)
}

// 下载文件接口
func DownloadFile(c *gin.Context) {
	taskUUID := c.Query("taskId")

	if taskUUID == "" {
		response.Failed(c, http.StatusBadRequest, "任务ID不能为空")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	// 查找下载任务
	var task models.DownloadTask
	if err := db.Where("uuid = ?", taskUUID).First(&task).Error; err != nil {
		response.Failed(c, http.StatusNotFound, "下载任务不存在")
		return
	}

	// 检查任务是否已完成
	if task.Status == models.DownloadTaskStatusCompleted {
		response.Failed(c, http.StatusBadRequest, "该任务已完成下载，不允许再次下载")
		return
	}

	// 检查文件是否存在
	fileInfo, err := os.Stat(task.FilePath)
	if os.IsNotExist(err) {
		response.Failed(c, http.StatusNotFound, "文件不存在")
		return
	}

	// 更新任务状态为下载中
	task.Status = models.DownloadTaskStatusDownloading
	db.Save(&task)

	// 处理Range请求
	rangeHeader := c.GetHeader("Range")
	var start, end int64 = 0, fileInfo.Size() - 1

	if rangeHeader != "" {
		_, err := fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
		if err != nil {
			start = 0
			end = fileInfo.Size() - 1
		}
	}

	// 打开文件
	file, err := os.Open(task.FilePath)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "打开文件失败")
		return
	}
	defer file.Close()

	// 设置响应头
	fileName := filepath.Base(task.FilePath)
	if rangeHeader != "" {
		c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileInfo.Size()))
		c.Header("Content-Length", fmt.Sprintf("%d", end-start+1))
		c.Status(http.StatusPartialContent)
	} else {
		c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	}

	// 定位到起始位置
	_, err = file.Seek(start, 0)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "文件定位失败")
		return
	}

	// 流式传输文件
	writer := c.Writer
	done := make(chan bool)

	// 异步更新下载进度
	go func() {
		defer close(done)

		// 模拟进度更新
		var downloaded int64 = 0
		buffer := make([]byte, 8192)
		for {
			n, err := file.Read(buffer)
			if n > 0 {
				writer.Write(buffer[:n])
				downloaded += int64(n)

				// 每读取一些数据就更新进度
				if downloaded%102400 == 0 {
					progress := float64(downloaded) / float64(fileInfo.Size()) * 100
					task.Progress = progress
					task.DownloadedBytes = downloaded
					db.Save(&task)
				}
			}
			if err != nil {
				break
			}
		}

		// 下载完成，更新任务状态
		task.Status = models.DownloadTaskStatusCompleted
		task.Progress = 100
		task.DownloadedBytes = fileInfo.Size()
		now := time.Now()
		task.CompletedAt = &now
		db.Save(&task)
	}()

	// 等待下载完成
	<-done
}

// 旧的下载接口（保持兼容）
func DownloadPluginFile(c *gin.Context) {
	filename := c.Param("filename")
	license := c.Query("license")
	taskUUID := c.Query("taskId")

	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	if license == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "License is required"})
		return
	}

	// 连接数据库
	db, err := database.GetDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}

	// 查找对应的插件
	var plugin models.Plugin
	if err := db.Where("file_path LIKE ?", "./uploads/plugins/"+filename).First(&plugin).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plugin not found"})
		return
	}

	// 验证插件状态
	if plugin.Status != models.PluginStatusPublished {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Plugin not published"})
		return
	}

	// 验证许可证
	var licenseObj models.License
	if err := db.Where("serial_number = ? AND status = ?", license, models.LicenseStatusApproved).First(&licenseObj).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid license"})
		return
	}

	// 验证许可证是否过期
	if licenseObj.ExpiryDate.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "License expired"})
		return
	}

	// 验证许可证类型是否与插件类型匹配
	if licenseObj.LicenseType != plugin.LicenseType {
		c.JSON(http.StatusBadRequest, gin.H{"error": "License type does not match plugin type"})
		return
	}

	// 验证许可证产品是否与插件产品匹配
	if licenseObj.ProductUUID != plugin.ProductUUID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "License product does not match plugin product"})
		return
	}

	pluginDir := "./uploads/plugins/"
	filePath := pluginDir + filename

	if !strings.HasPrefix(filepath.Clean(filePath), filepath.Join(pluginDir)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to access file"})
		return
	}

	// 如果没有taskId，创建新的下载任务
	if taskUUID == "" {
		taskUUID = uuid.New().String()
		clientIP := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")

		task := models.DownloadTask{
			UUID:       taskUUID,
			PluginUUID: plugin.UUID,
			License:    license,
			FilePath:   filePath,
			FileSize:   fileInfo.Size(),
			Status:     models.DownloadTaskStatusPending,
			IP:         clientIP,
			UserAgent:  userAgent,
			StartedAt:  time.Now(),
		}

		if err := db.Create(&task).Error; err != nil {
			utils.Logger.Errorw("DownloadPluginFile create task failed", "error", err)
		}

		// 返回下载任务ID
		c.JSON(http.StatusOK, gin.H{
			"taskId":      taskUUID,
			"downloadUrl": fmt.Sprintf("http://localhost:8082/download/%s?license=%s&taskId=%s", filename, license, taskUUID),
			"fileSize":    fileInfo.Size(),
		})
		return
	}

	// 查找下载任务
	var task models.DownloadTask
	if err := db.Where("uuid = ?", taskUUID).First(&task).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Download task not found"})
		return
	}

	// 检查任务是否已完成
	if task.Status == models.DownloadTaskStatusCompleted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task already completed"})
		return
	}

	// 更新任务状态为下载中
	task.Status = models.DownloadTaskStatusDownloading
	db.Save(&task)

	// 处理Range请求
	rangeHeader := c.GetHeader("Range")
	var start, end int64 = 0, fileInfo.Size() - 1

	if rangeHeader != "" {
		_, err := fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
		if err != nil {
			start = 0
			end = fileInfo.Size() - 1
		}
	}

	// 更新下载进度
	downloadedBytes := end + 1
	progress := float64(downloadedBytes) / float64(fileInfo.Size()) * 100

	task.DownloadedBytes = downloadedBytes
	task.Progress = progress
	if downloadedBytes >= fileInfo.Size() {
		task.Status = models.DownloadTaskStatusCompleted
		now := time.Now()
		task.CompletedAt = &now
	}
	db.Save(&task)

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

	// 设置响应头
	if rangeHeader != "" {
		c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileInfo.Size()))
		c.Header("Content-Length", fmt.Sprintf("%d", end-start+1))
		c.Status(http.StatusPartialContent)
	} else {
		c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	}

	// 定位到起始位置
	_, err = file.Seek(start, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to seek file"})
		return
	}

	// 流式传输文件
	http.ServeContent(c.Writer, c.Request, filename, fileInfo.ModTime(), file)
}

func GetDownloadProgress(c *gin.Context) {
	taskUUID := c.Param("taskId")

	if taskUUID == "" {
		response.Failed(c, http.StatusBadRequest, "Task ID is required")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "Database connection failed")
		return
	}

	var task models.DownloadTask
	if err := db.Preload("Plugin").Where("uuid = ?", taskUUID).First(&task).Error; err != nil {
		response.Failed(c, http.StatusNotFound, "Download task not found")
		return
	}

	response.Success(c, task, 1)
}
