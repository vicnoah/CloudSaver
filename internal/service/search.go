package service

import (
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// SearchService 搜索服务
type SearchService struct {
	httpClient *http.Client
}

// NewSearchService 创建搜索服务
func NewSearchService() *SearchService {
	return &SearchService{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SearchResult 搜索结果
type SearchResult struct {
	ID          string      `json:"id"`
	ChannelInfo ChannelInfo `json:"channelInfo"`
	List        []Resource  `json:"list"`
}

// ChannelInfo 频道信息
type ChannelInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ChannelLogo string `json:"channelLogo"`
}

// Resource 资源信息
type Resource struct {
	MessageID  string   `json:"messageId"`
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	PubDate    string   `json:"pubDate"`
	Image      string   `json:"image"`
	CloudLinks []string `json:"cloudLinks"`
	CloudType  string   `json:"cloudType"`
	Tags       []string `json:"tags"`
	Channel    string   `json:"channel"`
	ChannelID  string   `json:"channelId"`
}

// Search 搜索资源
func (s *SearchService) Search(keyword string, channelID string, messageID string) ([]SearchResult, error) {
	// TODO: 实现实际的搜索逻辑
	// 这里返回空结果，实际应该爬取 Telegram 频道
	return []SearchResult{}, nil
}

// extractCloudLinks 提取云盘链接
func (s *SearchService) extractCloudLinks(text string) ([]string, string) {
	links := []string{}
	cloudType := ""

	patterns := map[string]*regexp.Regexp{
		"pan115": regexp.MustCompile(`https://115\.com/s/[a-zA-Z0-9]+\?password=[a-zA-Z0-9]+`),
		"quark":  regexp.MustCompile(`https://pan\.quark\.cn/s/[a-zA-Z0-9]+`),
		"baidu":  regexp.MustCompile(`https://pan\.baidu\.com/s/[a-zA-Z0-9-_]+`),
		"aliyun": regexp.MustCompile(`https://www\.alipan\.com/s/[a-zA-Z0-9]+`),
	}

	for typ, pattern := range patterns {
		matches := pattern.FindAllString(text, -1)
		if len(matches) > 0 {
			links = append(links, matches...)
			if cloudType == "" {
				cloudType = typ
			}
		}
	}

	return links, cloudType
}

// fetchTelegramChannel 获取 Telegram 频道内容
func (s *SearchService) fetchTelegramChannel(channelID string, keyword string) (*SearchResult, error) {
	baseURL := "https://t.me/s/" + channelID
	if keyword != "" {
		baseURL += "?q=" + url.QueryEscape(keyword)
	}

	resp, err := s.httpClient.Get(baseURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 简单的HTML解析（实际应使用 goquery 等库）
	html := string(body)
	
	// 提取频道logo
	logoPattern := regexp.MustCompile(`<img[^>]*class="tgme_page_photo_image"[^>]*src="([^"]+)"`)
	logoMatches := logoPattern.FindStringSubmatch(html)
	channelLogo := ""
	if len(logoMatches) > 1 {
		channelLogo = logoMatches[1]
	}

	result := &SearchResult{
		ID: channelID,
		ChannelInfo: ChannelInfo{
			ID:          channelID,
			Name:        channelID,
			ChannelLogo: channelLogo,
		},
		List: []Resource{},
	}

	// 简单解析消息（实际应使用更完善的HTML解析器）
	messagePattern := regexp.MustCompile(`(?s)class="tgme_widget_message[^"]*".*?data-post="[^/]+/(\d+)"`)
	matches := messagePattern.FindAllStringSubmatch(html, -1)
	
	for _, match := range matches {
		if len(match) > 1 {
			resource := Resource{
				MessageID: match[1],
				Channel:   channelID,
				ChannelID: channelID,
			}
			
			// 提取云盘链接
			links, cloudType := s.extractCloudLinks(html)
			if len(links) > 0 {
				resource.CloudLinks = links
				resource.CloudType = cloudType
				result.List = append(result.List, resource)
			}
		}
	}

	return result, nil
}

// unique 去重
func unique(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			if entry != "" {
				list = append(list, entry)
			}
		}
	}
	return list
}

// contains 检查是否包含
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.Contains(strings.ToLower(s), strings.ToLower(item)) {
			return true
		}
	}
	return false
}
