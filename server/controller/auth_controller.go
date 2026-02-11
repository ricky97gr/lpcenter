package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ricky97gr/lpcenter/server/database"
	"github.com/ricky97gr/lpcenter/server/models"
	"github.com/ricky97gr/lpcenter/server/response"
	"github.com/ricky97gr/lpcenter/server/utils"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Logger.Errorw("Login parameter error", "error", err)
		response.Failed(c, http.StatusBadRequest, "参数错误")
		return
	}

	db, err := database.GetDB()
	if err != nil {
		utils.Logger.Errorw("Login database connection failed", "error", err)
		response.Failed(c, http.StatusInternalServerError, "数据库连接失败")
		return
	}

	var user models.User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.Logger.Errorf("Login user not found, username: %s, error: %v", req.Username, err)
		response.Failed(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	if user.Status != models.UserStatusActive {
		utils.Logger.Warnf("Login user disabled, username: %s, status: %d", req.Username, user.Status)
		response.Failed(c, http.StatusForbidden, "用户已被禁用")
		return
	}

	if req.Password != user.Password {
		utils.Logger.Errorf("Login password incorrect, username: %s", req.Username)
		response.Failed(c, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	token, err := utils.GenerateToken(user.UUID, user.Username)
	if err != nil {
		utils.Logger.Errorf("Login generate token failed, username: %s, error: %v", req.Username, err)
		response.Failed(c, http.StatusInternalServerError, "生成令牌失败")
		return
	}

	utils.Logger.Infof("Login success, userId: %s, username: %s", user.UUID, user.Username)
	response.Success(c, LoginResponse{
		Token: token,
		User:  toUserResponse(user),
	}, 1)
}

func Logout(c *gin.Context) {
	utils.Logger.Infow("Logout success")
	response.Success(c, nil, 0)
}
