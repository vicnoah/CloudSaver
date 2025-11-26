package handler

import (
	"net/http"

	"cloudsaver/internal/api/response"
	"cloudsaver/internal/service"

	"github.com/gin-gonic/gin"
)

// SearchHandler 搜索处理器
type SearchHandler struct {
	searchService *service.SearchService
}

// NewSearchHandler 创建搜索处理器
func NewSearchHandler(searchService *service.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// Search 搜索资源
func (h *SearchHandler) Search(c *gin.Context) {
	keyword := c.Query("keyword")
	channelID := c.Query("channelId")
	messageID := c.Query("messageId")

	result, err := h.searchService.Search(keyword, channelID, messageID)
	if err != nil {
		c.JSON(http.StatusOK, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(result))
}
