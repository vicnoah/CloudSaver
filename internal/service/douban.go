package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// DoubanService 豆瓣服务
type DoubanService struct {
	httpClient *http.Client
	baseURL    string
}

// NewDoubanService 创建豆瓣服务
func NewDoubanService() *DoubanService {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	return &DoubanService{
		httpClient: client,
		baseURL:    "https://movie.douban.com/j",
	}
}

// DoubanItem 豆瓣条目
type DoubanItem struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Rate  string `json:"rate"`
	Cover string `json:"cover"`
	URL   string `json:"url"`
	IsNew bool   `json:"is_new"`
}

// DoubanResponse 豆瓣API响应
type doubanResponse struct {
	Subjects []DoubanItem `json:"subjects"`
}

// GetHotList 获取豆瓣热门列表
func (s *DoubanService) GetHotList(typ, tag, pageLimit, pageStart string) ([]DoubanItem, error) {
	// 构建请求URL
	params := url.Values{}
	params.Add("type", typ)
	params.Add("tag", tag)
	params.Add("page_limit", pageLimit)
	params.Add("page_start", pageStart)

	reqURL := fmt.Sprintf("%s/search_subjects?%s", s.baseURL, params.Encode())

	// 创建请求
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://movie.douban.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	// 发送请求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析JSON
	var result doubanResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Subjects, nil
}
