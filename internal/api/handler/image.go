package handler

import (
	"io"
	"net/http"
	"time"

	"cloudsaver/internal/api/response"

	"github.com/gin-gonic/gin"
)

// ImageHandler 图片代理处理器
type ImageHandler struct {
	httpClient *http.Client
}

// NewImageHandler 创建图片代理处理器
func NewImageHandler() *ImageHandler {
	return &ImageHandler{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ProxyImage 代理图片请求
func (h *ImageHandler) ProxyImage(c *gin.Context) {
	imageURL := c.Query("url")
	if imageURL == "" {
		c.JSON(http.StatusBadRequest, response.Error("缺少图片URL参数"))
		return
	}

	// 创建请求
	req, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("创建请求失败"))
		return
	}

	// 设置请求头，伪装成浏览器访问，避免防盗链
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Referer", "https://movie.douban.com/")
	req.Header.Set("Accept", "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	// 发送请求
	resp, err := h.httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("获取图片失败"))
		return
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, response.Error("图片服务器返回错误"))
		return
	}

	// 设置响应头
	c.Header("Content-Type", resp.Header.Get("Content-Type"))
	c.Header("Cache-Control", "public, max-age=31536000") // 缓存一年
	
	// 复制图片数据到响应
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		// 如果已经开始写入响应，无法返回错误JSON
		return
	}
}
