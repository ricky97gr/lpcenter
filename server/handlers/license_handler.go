package handlers

import (
	"lpcenter/database"
	"lpcenter/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateLicenseRequest(c *gin.Context) {
	var license models.License
	if err := c.ShouldBindJSON(&license); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	license.Status = "pending"

	if err := database.GetDB().Create(&license).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create license request"})
		return
	}

	c.JSON(http.StatusCreated, license)
}

func GetLicenseRequests(c *gin.Context) {
	var licenses []models.License
	if err := database.GetDB().Find(&licenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get license requests"})
		return
	}

	c.JSON(http.StatusOK, licenses)
}

func GetLicenseRequest(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid license ID"})
		return
	}

	var license models.License
	if err := database.GetDB().First(&license, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "License request not found"})
		return
	}

	c.JSON(http.StatusOK, license)
}

func ApproveLicense(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid license ID"})
		return
	}

	var license models.License
	if err := database.GetDB().First(&license, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "License request not found"})
		return
	}

	license.Status = "approved"
	if err := database.GetDB().Save(&license).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve license"})
		return
	}

	c.JSON(http.StatusOK, license)
}

func RejectLicense(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid license ID"})
		return
	}

	var license models.License
	if err := database.GetDB().First(&license, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "License request not found"})
		return
	}

	license.Status = "rejected"
	if err := database.GetDB().Save(&license).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject license"})
		return
	}

	c.JSON(http.StatusOK, license)
}

func GetUserLicenses(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var licenses []models.License
	if err := database.GetDB().Where("user_id = ?", userID).Find(&licenses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user licenses"})
		return
	}

	c.JSON(http.StatusOK, licenses)
}
