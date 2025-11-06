package handler

import (
	"net/http"

	"cloudsaver/internal/api/request"
	"cloudsaver/internal/api/response"
	"cloudsaver/internal/service"

	"github.com/gin-gonic/gin"
)

// SettingHandler 设置处理器
type SettingHandler struct {
	settingService *service.SettingService
}

// NewSettingHandler 创建设置处理器
func NewSettingHandler(settingService *service.SettingService) *SettingHandler {
	return &SettingHandler{
		settingService: settingService,
	}
}

// Get 获取用户设置
func (h *SettingHandler) Get(c *gin.Context) {
	userUUID := c.GetString("user_uuid")
	if userUUID == "" {
		c.JSON(http.StatusUnauthorized, response.Error("未授权"))
		return
	}

	setting, err := h.settingService.GetUserSetting(userUUID)
	if err != nil {
		c.JSON(http.StatusOK, response.Error("获取设置失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(setting))
}

// Save 保存用户设置
func (h *SettingHandler) Save(c *gin.Context) {
	userUUID := c.GetString("user_uuid")
	if userUUID == "" {
		c.JSON(http.StatusUnauthorized, response.Error("未授权"))
		return
	}

	var req request.SaveSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误: "+err.Error()))
		return
	}

	err := h.settingService.SaveUserSetting(userUUID, req.Cloud115Cookie, req.QuarkCookie)
	if err != nil {
		c.JSON(http.StatusOK, response.Error("保存设置失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithMessage(nil, "保存成功"))
}
