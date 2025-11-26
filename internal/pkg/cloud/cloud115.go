package cloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"cloudsaver/internal/api/response"
	"cloudsaver/internal/repository"
)

// Cloud115ListItem 115网盘列表项
type Cloud115ListItem struct {
	CID  string `json:"cid"`
	N    string `json:"n"`
	S    int64  `json:"s"`
	NS   int    `json:"ns"`
}

// Cloud115Response 115网盘API响应
type Cloud115Response struct {
	State bool        `json:"state"`
	Error string      `json:"error"`
	Data  interface{} `json:"data"`
}

// Cloud115FilesResponse 115文件列表响应
type Cloud115FilesResponse struct {
	State bool                `json:"state"`
	Error string              `json:"error"`
	Data  []Cloud115ListItem  `json:"data"`
	Path  []Cloud115PathItem  `json:"path"`
}

// Cloud115PathItem 路径项
type Cloud115PathItem struct {
	CID  string `json:"cid"`
	Name string `json:"name"`
}

// Cloud115ShareResponse 115分享响应
type Cloud115ShareResponse struct {
	State bool `json:"state"`
	Error string `json:"error"`
	Data  struct {
		List []Cloud115ListItem `json:"list"`
	} `json:"data"`
}

// Cloud115Service 115网盘服务
type Cloud115Service struct {
	settingRepo *repository.SettingRepository
	cookie      string
	client      *http.Client
}

// NewCloud115Service 创建115网盘服务
func NewCloud115Service(settingRepo *repository.SettingRepository) *Cloud115Service {
	return &Cloud115Service{
		settingRepo: settingRepo,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetCookie 设置Cookie
func (s *Cloud115Service) SetCookie(userUUID string) error {
	setting, err := s.settingRepo.GetUserSetting(userUUID)
	if err != nil {
		return err
	}

	if setting.Cloud115Cookie == "" {
		return fmt.Errorf("请先设置115网盘Cookie")
	}

	s.cookie = setting.Cloud115Cookie
	return nil
}

// GetShareInfo 获取分享信息
func (s *Cloud115Service) GetShareInfo(ctx context.Context, shareCode string, passcode string) (*response.ShareInfoResponse, error) {
	params := url.Values{}
	params.Set("share_code", shareCode)
	params.Set("receive_code", passcode)
	params.Set("offset", "0")
	params.Set("limit", "20")
	params.Set("cid", "")

	apiURL := "https://webapi.115.com/share/snap?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	s.setHeaders(req)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Cloud115ShareResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if !result.State || len(result.Data.List) == 0 {
		return nil, fmt.Errorf("未找到文件信息: %s", result.Error)
	}

	fileItems := make([]response.FileItem, 0, len(result.Data.List))
	var totalSize int64
	for _, item := range result.Data.List {
		fileItems = append(fileItems, response.FileItem{
			FileID:   item.CID,
			FileName: item.N,
			FileType: item.NS,
		})
		totalSize += item.S
	}

	return &response.ShareInfoResponse{
		List:     fileItems,
		FileSize: totalSize,
	}, nil
}

// GetFolderList 获取文件夹列表
func (s *Cloud115Service) GetFolderList(ctx context.Context, parentID string) (*response.FolderListResponse, error) {
	if parentID == "" {
		parentID = "0"
	}

	params := url.Values{}
	params.Set("aid", "1")
	params.Set("cid", parentID)
	params.Set("o", "user_ptime")
	params.Set("asc", "1")
	params.Set("offset", "0")
	params.Set("show_dir", "1")
	params.Set("limit", "50")
	params.Set("type", "0")
	params.Set("format", "json")
	params.Set("star", "0")
	params.Set("suffix", "")
	params.Set("natsort", "0")
	params.Set("snap", "0")
	params.Set("record_open_time", "1")
	params.Set("fc_mix", "0")

	apiURL := "https://webapi.115.com/files?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	s.setHeaders(req)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Cloud115FilesResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if !result.State {
		return nil, fmt.Errorf("获取115pan目录列表失败: %s", result.Error)
	}

	folders := make([]response.FolderItem, 0)
	for _, item := range result.Data {
		if item.CID != "" && item.NS > 0 {
			pathItems := make([]response.FolderItem, 0, len(result.Path))
			for _, p := range result.Path {
				pathItems = append(pathItems, response.FolderItem{
					CID:  p.CID,
					Name: p.Name,
				})
			}
			folders = append(folders, response.FolderItem{
				CID:  item.CID,
				Name: item.N,
				Path: pathItems,
			})
		}
	}

	return &response.FolderListResponse{
		Folders: folders,
	}, nil
}

// SaveSharedFile 保存分享文件
func (s *Cloud115Service) SaveSharedFile(ctx context.Context, shareCode, passcode, folderID string, fileIDs []string, fileTokens []string, pwdID, stoken string) error {
	if len(fileIDs) == 0 {
		return fmt.Errorf("文件ID不能为空")
	}

	data := url.Values{}
	data.Set("cid", folderID)
	data.Set("share_code", shareCode)
	data.Set("receive_code", passcode)
	data.Set("file_id", fileIDs[0])

	req, err := http.NewRequestWithContext(ctx, "POST", "https://webapi.115.com/share/receive", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.setHeaders(req)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result Cloud115Response
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if !result.State {
		return fmt.Errorf("保存115pan文件失败: %s", result.Error)
	}

	return nil
}

// setHeaders 设置请求头
func (s *Cloud115Service) setHeaders(req *http.Request) {
	req.Header.Set("Host", "webapi.115.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("xweb_xhr", "1")
	req.Header.Set("Origin", "")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 MicroMessenger/6.8.0(0x16080000) NetType/WIFI MiniProgramEnv/Mac MacWechat/WMPF MacWechat/3.8.9(0x13080910) XWEB/1227")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://servicewechat.com/wx2c744c010a61b0fa/94/page-frame.html")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cookie", s.cookie)
}
