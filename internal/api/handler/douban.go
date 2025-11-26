package handler

import (
	"net/http"

	"cloudsaver/internal/api/response"
	"cloudsaver/internal/service"

	"github.com/gin-gonic/gin"
)

// DoubanHandler 豆瓣处理器
type DoubanHandler struct {
	doubanService *service.DoubanService
}

// NewDoubanHandler 创建豆瓣处理器
func NewDoubanHandler(doubanService *service.DoubanService) *DoubanHandler {
	return &DoubanHandler{
		doubanService: doubanService,
	}
}

// GetHotList 获取豆瓣热门列表
func (h *DoubanHandler) GetHotList(c *gin.Context) {
	typ := c.Query("type")
	tag := c.Query("tag")
	pageLimit := c.DefaultQuery("page_limit", "20")
	pageStart := c.DefaultQuery("page_start", "0")

	result, err := h.doubanService.GetHotList(typ, tag, pageLimit, pageStart)
	if err != nil {
		c.JSON(http.StatusOK, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(result))
}
