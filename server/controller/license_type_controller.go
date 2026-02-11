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

type LicenseTypeRequest struct {
	Name      string `json:"name" binding:"required"`
	Code      string `json:"code"`
	IsPaid    bool   `json:"isPaid"`
	CreatedBy string `json:"createdBy"`
}

func CreateLicenseType(c *gin.Context) {
	var req LicenseTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Errorw("CreateLicenseType parameter error", "error", err)
		response.Failed(c, http.StatusBadRequest, "参数错误")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("CreateLicenseType database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	code := req.Code
	if code == "" {
		code = strings.ToLower(req.Name)
		code = strings.ReplaceAll(code, "版", "")
		code = strings.ReplaceAll(code, " ", "")
	}

	licenseType := models.LicenseType{
		UUID:      uuid.New().String(),
		Name:      req.Name,
		Code:      code,
		IsPaid:    req.IsPaid,
		CreatedBy: req.CreatedBy,
	}

	if err := db.Create(&licenseType).Error; err != nil {
		utils.Logger.Errorw("CreateLicenseType failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "创建授权类型失败")
		return
	}

	utils.Logger.Infow("CreateLicenseType success", "licenseTypeId", licenseType.UUID, "name", licenseType.Name)
	response.Success(c, licenseType, 1)
}

func GetAllLicenseTypes(c *gin.Context) {
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("GetAllLicenseTypes database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	pageQuery, err := pagination.GetPageQuery(c)
	if err != nil {
		utils.Logger.Errorw("GetAllLicenseTypes page query error", "error", err)
		response.Failed(c, http.StatusBadRequest, "分页参数错误")
		return
	}

	var licenseTypes []models.LicenseType
	var total int64

	countDB := db.Model(&models.LicenseType{})
	for _, cond := range pageQuery.Conditions {
		countDB = countDB.Scopes(pagination.QueryFilter(cond.Field, cond.Value, cond.Operation))
	}
	if err := countDB.Count(&total).Error; err != nil {
		utils.Logger.Errorw("GetAllLicenseTypes count failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "获取授权类型列表失败")
		return
	}

	queryDB := db.Scopes(pagination.ParseQuery(pageQuery))
	if err := queryDB.Find(&licenseTypes).Error; err != nil {
		utils.Logger.Errorw("GetAllLicenseTypes find failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "获取授权类型列表失败")
		return
	}

	utils.Logger.Infow("GetAllLicenseTypes success", "count", total, "page", pageQuery.Page, "pageSize", pageQuery.PageSize)
	response.Success(c, licenseTypes, total)
}

func GetLicenseType(c *gin.Context) {
	id := c.Param("id")
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("GetLicenseType database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var licenseType models.LicenseType

	if err := db.First(&licenseType, "uuid = ?", id).Error; err != nil {
		utils.Logger.Errorw("GetLicenseType not found", "id", id, "error", err)
		response.Failed(c, http.StatusNotFound, "授权类型不存在")
		return
	}

	utils.Logger.Infow("GetLicenseType success", "licenseTypeId", licenseType.UUID)
	response.Success(c, licenseType, 1)
}

func UpdateLicenseType(c *gin.Context) {
	id := c.Param("id")
	var req LicenseTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Errorw("UpdateLicenseType parameter error", "id", id, "error", err)
		response.Failed(c, http.StatusBadRequest, "参数错误")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("UpdateLicenseType database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var licenseType models.LicenseType

	if err := db.First(&licenseType, "uuid = ?", id).Error; err != nil {
		utils.Logger.Errorw("UpdateLicenseType not found", "id", id, "error", err)
		response.Failed(c, http.StatusNotFound, "授权类型不存在")
		return
	}

	code := req.Code
	if code == "" {
		code = strings.ToLower(req.Name)
		code = strings.ReplaceAll(code, "版", "")
		code = strings.ReplaceAll(code, " ", "")
	}

	licenseType.Name = req.Name
	licenseType.Code = code
	licenseType.IsPaid = req.IsPaid
	licenseType.CreatedBy = req.CreatedBy

	if err := db.Save(&licenseType).Error; err != nil {
		utils.Logger.Errorw("UpdateLicenseType failed", "id", id, "error", err)
		response.Failed(c, http.StatusInternalServerError, "更新授权类型失败")
		return
	}

	utils.Logger.Infow("UpdateLicenseType success", "licenseTypeId", licenseType.UUID)
	response.Success(c, licenseType, 1)
}

func DeleteLicenseType(c *gin.Context) {
	id := c.Param("id")
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("DeleteLicenseType database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	if err := db.Where("uuid = ?", id).Delete(&models.LicenseType{}).Error; err != nil {
		utils.Logger.Errorw("DeleteLicenseType failed", "id", id, "error", err)
		response.Failed(c, http.StatusInternalServerError, "删除授权类型失败")
		return
	}

	utils.Logger.Infow("DeleteLicenseType success", "licenseTypeId", id)
	response.Success(c, nil, 0)
}
