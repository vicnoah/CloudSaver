package handler

import (
	"net/http"

	"cloudsaver/internal/api/request"
	"cloudsaver/internal/api/response"
	"cloudsaver/internal/pkg/cloud"

	"github.com/gin-gonic/gin"
)

// QuarkHandler 夸克网盘处理器
type QuarkHandler struct {
	quarkService *cloud.QuarkService
}

// NewQuarkHandler 创建夸克网盘处理器
func NewQuarkHandler(quarkService *cloud.QuarkService) *QuarkHandler {
	return &QuarkHandler{
		quarkService: quarkService,
	}
}

// GetShareInfo 获取分享信息
func (h *QuarkHandler) GetShareInfo(c *gin.Context) {
	var req request.GetShareInfoRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误"))
		return
	}

	result, err := h.quarkService.GetShareInfo(c.Request.Context(), req.ShareCode, req.Passcode)
	if err != nil {
		c.JSON(http.StatusOK, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(result))
}

// GetFolderList 获取文件夹列表
func (h *QuarkHandler) GetFolderList(c *gin.Context) {
	userUUID := c.GetString("user_uuid")
	if userUUID == "" {
		c.JSON(http.StatusUnauthorized, response.Error("未授权"))
		return
	}

	// 设置Cookie
	if err := h.quarkService.SetCookie(userUUID); err != nil {
		c.JSON(http.StatusOK, response.Error(err.Error()))
		return
	}

	var req request.GetFolderListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误"))
		return
	}

	result, err := h.quarkService.GetFolderList(c.Request.Context(), req.ParentCid)
	if err != nil {
		c.JSON(http.StatusOK, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(result))
}

// SaveFile 保存文件
func (h *QuarkHandler) SaveFile(c *gin.Context) {
	userUUID := c.GetString("user_uuid")
	if userUUID == "" {
		c.JSON(http.StatusUnauthorized, response.Error("未授权"))
		return
	}

	// 设置Cookie
	if err := h.quarkService.SetCookie(userUUID); err != nil {
		c.JSON(http.StatusOK, response.Error(err.Error()))
		return
	}

	var req request.SaveFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("请求参数错误"))
		return
	}

	err := h.quarkService.SaveSharedFile(c.Request.Context(), req.ShareCode, req.Passcode, req.FolderID, req.FileIDs, req.FileTokens, req.PwdID, req.Stoken)
	if err != nil {
		c.JSON(http.StatusOK, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessWithMessage(nil, "转存成功"))
}
