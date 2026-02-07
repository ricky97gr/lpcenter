package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ricky97gr/lpcenter/server/database"
	"github.com/ricky97gr/lpcenter/server/models"
	"github.com/ricky97gr/lpcenter/server/pagination"
	"github.com/ricky97gr/lpcenter/server/response"
	"github.com/ricky97gr/lpcenter/server/utils"
)

type LicenseRequest struct {
	SerialNumber string `json:"serialNumber" binding:"required"`
	ProductID    uint   `json:"productId" binding:"required"`
	Version      string `json:"version" binding:"required"`
	ExpiryDate   string `json:"expiryDate" binding:"required"`
	Remarks      string `json:"remarks"`
}

func CreateLicense(c *gin.Context) {
	var req LicenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Errorw("CreateLicense parameter error", "error", err)
		response.Failed(c, http.StatusBadRequest, "参数错误")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("CreateLicense database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var product models.Product
	if err := db.First(&product, req.ProductID).Error; err != nil {
		utils.Logger.Errorw("CreateLicense product not found", "productId", req.ProductID, "error", err)
		response.Failed(c, http.StatusNotFound, "产品不存在")
		return
	}

	license := models.License{
		UUID:         uuid.New().String(),
		SerialNumber: req.SerialNumber,
		ProductID:    req.ProductID,
		Version:      req.Version,
		ExpiryDate:   parseExpiryDate(req.ExpiryDate),
		Remarks:      req.Remarks,
		Status:       models.LicenseStatusPending,
	}

	if err := db.Create(&license).Error; err != nil {
		utils.Logger.Errorw("CreateLicense failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "创建授权申请失败")
		return
	}

	utils.Logger.Infow("CreateLicense success", "licenseId", license.ID, "serialNumber", license.SerialNumber)
	response.Success(c, license, 1)
}

func GetAllLicenses(c *gin.Context) {
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("GetAllLicenses database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var licenses []models.License
	var total int64

	pageQuery, err := pagination.GetPageQuery(c)
	if err != nil {
		utils.Logger.Errorw("GetAllLicenses page query error", "error", err)
		response.Failed(c, http.StatusBadRequest, "分页参数错误")
		return
	}

	countDB := db.Model(&models.License{})
	for _, cond := range pageQuery.Conditions {
		countDB = countDB.Scopes(pagination.QueryFilter(cond.Field, cond.Value, cond.Operation))
	}
	if err := countDB.Count(&total).Error; err != nil {
		utils.Logger.Errorw("GetAllLicenses count failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "获取授权列表失败")
		return
	}

	queryDB := db.Preload("Product").Scopes(pagination.ParseQuery(pageQuery))
	if err := queryDB.Find(&licenses).Error; err != nil {
		utils.Logger.Errorw("GetAllLicenses find failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "获取授权列表失败")
		return
	}

	utils.Logger.Infow("GetAllLicenses success", "count", total, "page", pageQuery.Page, "pageSize", pageQuery.PageSize)
	response.Success(c, licenses, total)
}

func GetLicense(c *gin.Context) {
	id := c.Param("id")
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("GetLicense database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var license models.License

	if err := db.Preload("Product").First(&license, id).Error; err != nil {
		utils.Logger.Errorw("GetLicense not found", "id", id, "error", err)
		response.Failed(c, http.StatusNotFound, "授权不存在")
		return
	}

	utils.Logger.Infow("GetLicense success", "licenseId", license.ID)
	response.Success(c, license, 1)
}

func UpdateLicense(c *gin.Context) {
	id := c.Param("id")
	var req LicenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Errorw("UpdateLicense parameter error", "id", id, "error", err)
		response.Failed(c, http.StatusBadRequest, "参数错误")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("UpdateLicense database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var license models.License

	if err := db.First(&license, id).Error; err != nil {
		utils.Logger.Errorw("UpdateLicense not found", "id", id, "error", err)
		response.Failed(c, http.StatusNotFound, "授权不存在")
		return
	}

	var product models.Product
	if err := db.First(&product, req.ProductID).Error; err != nil {
		utils.Logger.Errorw("UpdateLicense product not found", "productId", req.ProductID, "error", err)
		response.Failed(c, http.StatusNotFound, "产品不存在")
		return
	}

	license.SerialNumber = req.SerialNumber
	license.ProductID = req.ProductID
	license.Version = req.Version
	license.ExpiryDate = parseExpiryDate(req.ExpiryDate)
	license.Remarks = req.Remarks

	if err := db.Save(&license).Error; err != nil {
		utils.Logger.Errorw("UpdateLicense failed", "id", id, "error", err)
		response.Failed(c, http.StatusInternalServerError, "更新授权失败")
		return
	}

	utils.Logger.Infow("UpdateLicense success", "licenseId", license.ID)
	response.Success(c, license, 1)
}

func DeleteLicense(c *gin.Context) {
	id := c.Param("id")
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("DeleteLicense database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	if err := db.Delete(&models.License{}, id).Error; err != nil {
		utils.Logger.Errorw("DeleteLicense failed", "id", id, "error", err)
		response.Failed(c, http.StatusInternalServerError, "删除授权失败")
		return
	}

	utils.Logger.Infow("DeleteLicense success", "licenseId", id)
	response.Success(c, nil, 0)
}

func ApproveLicense(c *gin.Context) {
	id := c.Param("id")
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("ApproveLicense database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var license models.License
	if err := db.Preload("Product").First(&license, id).Error; err != nil {
		utils.Logger.Errorw("ApproveLicense not found", "id", id, "error", err)
		response.Failed(c, http.StatusNotFound, "授权不存在")
		return
	}

	if license.Product == nil {
		utils.Logger.Errorw("ApproveLicense product not found", "id", id, "productId", license.ProductID)
		response.Failed(c, http.StatusNotFound, "产品不存在")
		return
	}

	licenseData := map[string]interface{}{
		"serialNumber": license.SerialNumber,
		"productId":    license.ProductID,
		"productName":  license.Product.Name,
		"version":      license.Version,
		"expiryDate":   license.ExpiryDate.Format("2006-01-02"),
		"expiryTime":   license.ExpiryDate.Unix(),
		"createdAt":    license.CreatedAt.Unix(),
		"status":       "approved",
	}

	licenseString, err := utils.GenerateLicenseString(licenseData)
	if err != nil {
		utils.Logger.Errorw("ApproveLicense generate license string failed", "id", id, "error", err)
		response.Failed(c, http.StatusInternalServerError, "生成授权字符串失败")
		return
	}

	license.Status = models.LicenseStatusApproved
	license.LicenseString = licenseString
	if err := db.Save(&license).Error; err != nil {
		utils.Logger.Errorw("ApproveLicense failed", "id", id, "error", err)
		response.Failed(c, http.StatusInternalServerError, "批准授权失败")
		return
	}

	utils.Logger.Infow("ApproveLicense success", "licenseId", license.ID)
	response.Success(c, license, 1)
}

func DownloadLicenseFile(c *gin.Context) {
	id := c.Param("id")
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("DownloadLicenseFile database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var license models.License
	if err := db.First(&license, id).Error; err != nil {
		utils.Logger.Errorw("DownloadLicenseFile not found", "id", id, "error", err)
		response.Failed(c, http.StatusNotFound, "授权不存在")
		return
	}

	if license.Status != models.LicenseStatusApproved {
		utils.Logger.Errorw("DownloadLicenseFile not approved", "id", id, "status", license.Status)
		response.Failed(c, http.StatusBadRequest, "授权未批准")
		return
	}

	if license.LicenseString == "" {
		utils.Logger.Errorw("DownloadLicenseFile license string empty", "id", id)
		response.Failed(c, http.StatusBadRequest, "授权字符串不存在")
		return
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+license.SerialNumber+".lic")
	c.Data(http.StatusOK, "application/octet-stream", []byte(license.LicenseString))

	utils.Logger.Infow("DownloadLicenseFile success", "licenseId", license.ID, "serialNumber", license.SerialNumber)
}

func RejectLicense(c *gin.Context) {
	id := c.Param("id")
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("RejectLicense database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var license models.License
	if err := db.First(&license, id).Error; err != nil {
		utils.Logger.Errorw("RejectLicense not found", "id", id, "error", err)
		response.Failed(c, http.StatusNotFound, "授权不存在")
		return
	}

	license.Status = models.LicenseStatusRejected
	if err := db.Save(&license).Error; err != nil {
		utils.Logger.Errorw("RejectLicense failed", "id", id, "error", err)
		response.Failed(c, http.StatusInternalServerError, "拒绝授权失败")
		return
	}

	utils.Logger.Infow("RejectLicense success", "licenseId", license.ID)
	response.Success(c, license, 1)
}

func parseExpiryDate(dateStr string) time.Time {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		t, _ = time.Parse("2006-01-02T15:04:05Z", dateStr)
	}
	return t
}
