package handler

import (
	"net/http"

	"cloudsaver/internal/api/request"
	"cloudsaver/internal/api/response"
	"cloudsaver/internal/pkg/cloud"

	"github.com/gin-gonic/gin"
)

// Cloud115Handler 115网盘处理器
type Cloud115Handler struct {
	cloud115Service *cloud.Cloud115Service
}

// NewCloud115Handler 创建115网盘处理器
func NewCloud115Handler(cloud115Service *cloud.Cloud115Service) *Cloud115Handler {
	return &Cloud115Handler{
		cloud115Service: cloud115Service,
	}
}

// GetShareInfo 获取分享信息
func (h *Cloud115Handler) GetShareInfo(c *gin.Context) {
	var req request.GetShareInfoRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误"))
		return
	}

	result, err := h.cloud115Service.GetShareInfo(c.Request.Context(), req.ShareCode, req.Passcode)
	if err != nil {
		c.JSON(http.StatusOK, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(result))
}

// GetFolderList 获取文件夹列表
func (h *Cloud115Handler) GetFolderList(c *gin.Context) {
	userUUID := c.GetString("user_uuid")
	if userUUID == "" {
		c.JSON(http.StatusUnauthorized, response.Error("未授权"))
		return
	}

	// 设置Cookie
	if err := h.cloud115Service.SetCookie(userUUID); err != nil {
		c.JSON(http.StatusOK, response.Error(err.Error()))
		return
	}

	var req request.GetFolderListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误"))
		return
	}

	result, err := h.cloud115Service.GetFolderList(c.Request.Context(), req.ParentCid)
	if err != nil {
		c.JSON(http.StatusOK, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(result))
}

// SaveFile 保存文件
func (h *Cloud115Handler) SaveFile(c *gin.Context) {
	userUUID := c.GetString("user_uuid")
	if userUUID == "" {
		c.JSON(http.StatusUnauthorized, response.Error("未授权"))
		return
	}

	// 设置Cookie
	if err := h.cloud115Service.SetCookie(userUUID); err != nil {
		c.JSON(http.StatusOK, response.Error(err.Error()))
		return
	}

	var req request.SaveFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误"))
		return
	}

	err := h.cloud115Service.SaveSharedFile(c.Request.Context(), req.ShareCode, req.Passcode, req.FolderID, req.FileIDs, nil, "", "")
	if err != nil {
		c.JSON(http.StatusOK, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithMessage(nil, "转存成功"))
}
