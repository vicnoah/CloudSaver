package handler

import (
	"net/http"

	"cloudsaver/internal/api/request"
	"cloudsaver/internal/api/response"
	"cloudsaver/internal/service"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	result, err := h.userService.Register(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusOK, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithMessage(result, "注册成功"))
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	result, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusOK, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithMessage(result, "登录成功"))
}
