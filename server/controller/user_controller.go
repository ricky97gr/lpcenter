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

type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Status   string `json:"status"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	UUID      string `json:"uuid"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

func toUserResponse(user models.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		UUID:      user.UUID,
		Username:  user.Username,
		Email:     user.Email,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func GetAllUsers(c *gin.Context) {
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("GetAllUsers database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	pageQuery, err := pagination.GetPageQuery(c)
	if err != nil {
		utils.Logger.Errorw("GetAllUsers page query error", "error", err)
		response.Failed(c, http.StatusBadRequest, "分页参数错误")
		return
	}

	var users []models.User
	var total int64

	countDB := db.Model(&models.User{})
	for _, cond := range pageQuery.Conditions {
		countDB = countDB.Scopes(pagination.QueryFilter(cond.Field, cond.Value, cond.Operation))
	}
	if err := countDB.Count(&total).Error; err != nil {
		utils.Logger.Errorw("GetAllUsers count failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "获取用户列表失败")
		return
	}

	queryDB := db.Scopes(pagination.ParseQuery(pageQuery))
	if err := queryDB.Find(&users).Error; err != nil {
		utils.Logger.Errorw("GetAllUsers find failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "获取用户列表失败")
		return
	}

	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, toUserResponse(user))
	}

	utils.Logger.Infow("GetAllUsers success", "count", total, "page", pageQuery.Page, "pageSize", pageQuery.PageSize)
	response.Success(c, userResponses, total)
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("GetUser database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		utils.Logger.Errorw("GetUser not found", "id", id, "error", err)
		response.Failed(c, http.StatusNotFound, "用户不存在")
		return
	}

	utils.Logger.Infow("GetUser success", "userId", user.ID)
	response.Success(c, toUserResponse(user), 1)
}

func CreateUser(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Errorw("CreateUser parameter error", "error", err)
		response.Failed(c, http.StatusBadRequest, "参数错误")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("CreateUser database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var existingUser models.User
	if err := db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		utils.Logger.Errorw("CreateUser user already exists", "username", req.Username, "email", req.Email)
		response.Failed(c, http.StatusBadRequest, "用户名或邮箱已存在")
		return
	}

	user := models.User{
		UUID:     uuid.New().String(),
		Username: req.Username,
		Email:    req.Email,
		Password: "123456",
		Status:   req.Status,
	}

	if user.Status == "" {
		user.Status = models.UserStatusActive
	}

	if err := db.Create(&user).Error; err != nil {
		utils.Logger.Errorw("CreateUser failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "创建用户失败")
		return
	}

	utils.Logger.Infow("CreateUser success", "userId", user.ID, "username", user.Username)
	response.Success(c, toUserResponse(user), 1)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Errorw("UpdateUser parameter error", "id", id, "error", err)
		response.Failed(c, http.StatusBadRequest, "参数错误")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("UpdateUser database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		utils.Logger.Errorw("UpdateUser not found", "id", id, "error", err)
		response.Failed(c, http.StatusNotFound, "用户不存在")
		return
	}

	var existingUser models.User
	if err := db.Where("username = ? OR email = ? AND id != ?", req.Username, req.Email, id).First(&existingUser).Error; err == nil {
		utils.Logger.Errorw("UpdateUser user already exists", "username", req.Username, "email", req.Email)
		response.Failed(c, http.StatusBadRequest, "用户名或邮箱已存在")
		return
	}

	user.Username = req.Username
	user.Email = req.Email
	if req.Status != "" {
		user.Status = req.Status
	}

	if err := db.Save(&user).Error; err != nil {
		utils.Logger.Errorw("UpdateUser failed", "id", id, "error", err)
		response.Failed(c, http.StatusInternalServerError, "更新用户失败")
		return
	}

	utils.Logger.Infow("UpdateUser success", "userId", user.ID)
	response.Success(c, toUserResponse(user), 1)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("DeleteUser database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	if err := db.Delete(&models.User{}, id).Error; err != nil {
		utils.Logger.Errorw("DeleteUser failed", "id", id, "error", err)
		response.Failed(c, http.StatusInternalServerError, "删除用户失败")
		return
	}

	utils.Logger.Infow("DeleteUser success", "userId", id)
	response.Success(c, nil, 0)
}

type ResetPasswordRequest struct {
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

func ResetPassword(c *gin.Context) {
	id := c.Param("id")
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Errorw("ResetPassword parameter error", "id", id, "error", err)
		response.Failed(c, http.StatusBadRequest, "参数错误")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("ResetPassword database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		utils.Logger.Errorw("ResetPassword not found", "id", id, "error", err)
		response.Failed(c, http.StatusNotFound, "用户不存在")
		return
	}

	user.Password = req.NewPassword
	if err := db.Save(&user).Error; err != nil {
		utils.Logger.Errorw("ResetPassword failed", "id", id, "error", err)
		response.Failed(c, http.StatusInternalServerError, "重置密码失败")
		return
	}

	utils.Logger.Infow("ResetPassword success", "userId", id)
	response.Success(c, nil, 0)
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

func ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Errorw("ChangePassword parameter error", "error", err)
		response.Failed(c, http.StatusBadRequest, "参数错误")
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		response.Failed(c, http.StatusUnauthorized, "未授权")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("ChangePassword database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		utils.Logger.Errorw("ChangePassword not found", "userId", userID, "error", err)
		response.Failed(c, http.StatusNotFound, "用户不存在")
		return
	}

	if req.OldPassword != user.Password {
		utils.Logger.Errorw("ChangePassword old password incorrect", "userId", userID)
		response.Failed(c, http.StatusBadRequest, "旧密码错误")
		return
	}

	user.Password = req.NewPassword
	if err := db.Save(&user).Error; err != nil {
		utils.Logger.Errorw("ChangePassword failed", "userId", userID, "error", err)
		response.Failed(c, http.StatusInternalServerError, "修改密码失败")
		return
	}

	utils.Logger.Infow("ChangePassword success", "userId", userID)
	response.Success(c, nil, 0)
}

func GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Failed(c, http.StatusUnauthorized, "未授权")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("GetCurrentUser database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		utils.Logger.Errorw("GetCurrentUser not found", "userId", userID, "error", err)
		response.Failed(c, http.StatusNotFound, "用户不存在")
		return
	}

	utils.Logger.Infow("GetCurrentUser success", "userId", userID)
	response.Success(c, toUserResponse(user), 1)
}
