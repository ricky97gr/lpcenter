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

type ProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Code        string `json:"code" binding:"required"`
	Status      string `json:"status"`
}

func CreateProduct(c *gin.Context) {
	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Errorw("CreateProduct parameter error", "error", err)
		response.Failed(c, http.StatusBadRequest, "参数错误")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("CreateProduct database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	product := models.Product{
		UUID:        uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Code:        req.Code,
		Status:      req.Status,
	}

	if req.Status == "" {
		product.Status = models.ProductStatusActive
	}

	if err := db.Create(&product).Error; err != nil {
		utils.Logger.Errorw("CreateProduct failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "创建产品失败")
		return
	}

	utils.Logger.Infow("CreateProduct success", "productId", product.UUID, "name", product.Name)
	response.Success(c, product, 1)
}

func GetAllProducts(c *gin.Context) {
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("GetAllProducts database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	pageQuery, err := pagination.GetPageQuery(c)
	if err != nil {
		utils.Logger.Errorw("GetAllProducts page query error", "error", err)
		response.Failed(c, http.StatusBadRequest, "分页参数错误")
		return
	}

	var products []models.Product
	var total int64

	countDB := db.Model(&models.Product{})
	for _, cond := range pageQuery.Conditions {
		countDB = countDB.Scopes(pagination.QueryFilter(cond.Field, cond.Value, cond.Operation))
	}
	if err := countDB.Count(&total).Error; err != nil {
		utils.Logger.Errorw("GetAllProducts count failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "获取产品列表失败")
		return
	}

	queryDB := db.Scopes(pagination.ParseQuery(pageQuery))
	if err := queryDB.Find(&products).Error; err != nil {
		utils.Logger.Errorw("GetAllProducts find failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "获取产品列表失败")
		return
	}

	utils.Logger.Infow("GetAllProducts success", "count", total, "page", pageQuery.Page, "pageSize", pageQuery.PageSize)
	response.Success(c, products, total)
}

func GetProduct(c *gin.Context) {
	id := c.Param("id")
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("GetProduct database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var product models.Product

	if err := db.First(&product, "uuid = ?", id).Error; err != nil {
		utils.Logger.Errorw("GetProduct not found", "id", id, "error", err)
		response.Failed(c, http.StatusNotFound, "产品不存在")
		return
	}

	utils.Logger.Infow("GetProduct success", "productId", product.UUID)
	response.Success(c, product, 1)
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Errorw("UpdateProduct parameter error", "id", id, "error", err)
		response.Failed(c, http.StatusBadRequest, "参数错误")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("UpdateProduct database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var product models.Product

	if err := db.First(&product, "uuid = ?", id).Error; err != nil {
		utils.Logger.Errorw("UpdateProduct not found", "id", id, "error", err)
		response.Failed(c, http.StatusNotFound, "产品不存在")
		return
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Code = req.Code
	if req.Status != "" {
		product.Status = req.Status
	}

	if err := db.Save(&product).Error; err != nil {
		utils.Logger.Errorw("UpdateProduct failed", "id", id, "error", err)
		response.Failed(c, http.StatusInternalServerError, "更新产品失败")
		return
	}

	utils.Logger.Infow("UpdateProduct success", "productId", product.UUID)
	response.Success(c, product, 1)
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("DeleteProduct database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	if err := db.Where("uuid = ?", id).Delete(&models.Product{}).Error; err != nil {
		utils.Logger.Errorw("DeleteProduct failed", "id", id, "error", err)
		response.Failed(c, http.StatusInternalServerError, "删除产品失败")
		return
	}

	utils.Logger.Infow("DeleteProduct success", "productId", id)
	response.Success(c, nil, 0)
}
