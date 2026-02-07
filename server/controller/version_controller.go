package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ricky97gr/lpcenter/server/database"
	"github.com/ricky97gr/lpcenter/server/models"
	"github.com/ricky97gr/lpcenter/server/pagination"
	"github.com/ricky97gr/lpcenter/server/response"
	"github.com/ricky97gr/lpcenter/server/utils"
)

type VersionRequest struct {
	Name      string `json:"name" binding:"required"`
	IsPaid    bool   `json:"isPaid"`
	CreatedBy string `json:"createdBy"`
}

func CreateVersion(c *gin.Context) {
	var req VersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Errorw("CreateVersion parameter error", "error", err)
		response.Failed(c, http.StatusBadRequest, "参数错误")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("CreateVersion database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	code := strings.ToLower(req.Name)
	code = strings.ReplaceAll(code, "版", "")
	code = strings.ReplaceAll(code, " ", "")

	version := models.Version{
		UUID:      uuid.New().String(),
		Name:      req.Name,
		Code:      code,
		IsPaid:    req.IsPaid,
		CreatedBy: req.CreatedBy,
	}

	if err := db.Create(&version).Error; err != nil {
		utils.Logger.Errorw("CreateVersion failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "创建版本失败")
		return
	}

	utils.Logger.Infow("CreateVersion success", "versionId", version.ID, "name", version.Name)
	response.Success(c, version, 1)
}

func GetAllVersions(c *gin.Context) {
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("GetAllVersions database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	pageQuery, err := pagination.GetPageQuery(c)
	if err != nil {
		utils.Logger.Errorw("GetAllVersions page query error", "error", err)
		response.Failed(c, http.StatusBadRequest, "分页参数错误")
		return
	}

	var versions []models.Version
	var total int64

	countDB := db.Model(&models.Version{})
	for _, cond := range pageQuery.Conditions {
		countDB = countDB.Scopes(pagination.QueryFilter(cond.Field, cond.Value, cond.Operation))
	}
	if err := countDB.Count(&total).Error; err != nil {
		utils.Logger.Errorw("GetAllVersions count failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "获取版本列表失败")
		return
	}

	queryDB := db.Scopes(pagination.ParseQuery(pageQuery))
	if err := queryDB.Find(&versions).Error; err != nil {
		utils.Logger.Errorw("GetAllVersions find failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "获取版本列表失败")
		return
	}

	utils.Logger.Infow("GetAllVersions success", "count", total, "page", pageQuery.Page, "pageSize", pageQuery.PageSize)
	response.Success(c, versions, total)
}

func GetVersion(c *gin.Context) {
	id := c.Param("id")
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("GetVersion database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var version models.Version

	if err := db.First(&version, id).Error; err != nil {
		utils.Logger.Errorw("GetVersion not found", "id", id, "error", err)
		response.Failed(c, http.StatusNotFound, "版本不存在")
		return
	}

	utils.Logger.Infow("GetVersion success", "versionId", version.ID)
	response.Success(c, version, 1)
}

func UpdateVersion(c *gin.Context) {
	id := c.Param("id")
	var req VersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Errorw("UpdateVersion parameter error", "id", id, "error", err)
		response.Failed(c, http.StatusBadRequest, "参数错误")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("UpdateVersion database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var version models.Version

	if err := db.First(&version, id).Error; err != nil {
		utils.Logger.Errorw("UpdateVersion not found", "id", id, "error", err)
		response.Failed(c, http.StatusNotFound, "版本不存在")
		return
	}

	version.Name = req.Name
	version.IsPaid = req.IsPaid
	version.CreatedBy = req.CreatedBy

	if err := db.Save(&version).Error; err != nil {
		utils.Logger.Errorw("UpdateVersion failed", "id", id, "error", err)
		response.Failed(c, http.StatusInternalServerError, "更新版本失败")
		return
	}

	utils.Logger.Infow("UpdateVersion success", "versionId", version.ID)
	response.Success(c, version, 1)
}

func DeleteVersion(c *gin.Context) {
	id := c.Param("id")
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("DeleteVersion database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	if err := db.Delete(&models.Version{}, id).Error; err != nil {
		utils.Logger.Errorw("DeleteVersion failed", "id", id, "error", err)
		response.Failed(c, http.StatusInternalServerError, "删除版本失败")
		return
	}

	utils.Logger.Infow("DeleteVersion success", "versionId", id)
	response.Success(c, nil, 0)
}
