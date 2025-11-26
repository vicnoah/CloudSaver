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

// QuarkTokenResponse 夸兌Token响应
type QuarkTokenResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Stoken string `json:"stoken"`
	} `json:"data"`
}

// QuarkShareDetail 夸兌分享详情
type QuarkShareDetail struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct {
		List []struct {
			FID           string `json:"fid"`
			FileName      string `json:"file_name"`
			FileType      int    `json:"file_type"`
			ShareFidToken string `json:"share_fid_token"`
		} `json:"list"`
		Share struct {
			Size int64 `json:"size"`
		} `json:"share"`
	} `json:"data"`
}

// QuarkFolderResponse 夸兌文件夹响应
type QuarkFolderResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct {
		List []struct {
			FID      string `json:"fid"`
			FileName string `json:"file_name"`
			FileType int    `json:"file_type"`
		} `json:"list"`
	} `json:"data"`
}

// QuarkSaveResponse 夸兌保存响应
type QuarkSaveResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// QuarkService 夸克网盘服务
type QuarkService struct {
	settingRepo *repository.SettingRepository
	cookie      string
	client      *http.Client
}

// NewQuarkService 创建夸克网盘服务
func NewQuarkService(settingRepo *repository.SettingRepository) *QuarkService {
	return &QuarkService{
		settingRepo: settingRepo,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetCookie 设置Cookie
func (s *QuarkService) SetCookie(userUUID string) error {
	setting, err := s.settingRepo.GetUserSetting(userUUID)
	if err != nil {
		return err
	}

	if setting.QuarkCookie == "" {
		return fmt.Errorf("请先设置夸克网盘Cookie")
	}

	s.cookie = setting.QuarkCookie
	return nil
}

// GetShareInfo 获取分享信息
func (s *QuarkService) GetShareInfo(ctx context.Context, pwdID string, passcode string) (*response.ShareInfoResponse, error) {
	// 步骤1: 获取stoken
	now := time.Now().UnixMilli()
	tokenURL := fmt.Sprintf("https://drive-h.quark.cn/1/clouddrive/share/sharepage/token?pr=ucpro&fr=pc&uc_param_str=&__dt=994&__t=%d", now)

	tokenData := map[string]string{
		"pwd_id":   pwdID,
		"passcode": passcode,
	}
	tokenBody, _ := json.Marshal(tokenData)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, bytes.NewReader(tokenBody))
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

	var tokenResp QuarkTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	if tokenResp.Status != 200 || tokenResp.Data.Stoken == "" {
		return nil, fmt.Errorf("获取夸兌分享信息失败: %s", tokenResp.Message)
	}

	// 步骤2: 获取分享文件列表
	return s.getShareList(ctx, pwdID, tokenResp.Data.Stoken)
}

// getShareList 获取分享列表
func (s *QuarkService) getShareList(ctx context.Context, pwdID, stoken string) (*response.ShareInfoResponse, error) {
	now := time.Now().UnixMilli()
	params := url.Values{}
	params.Set("pr", "ucpro")
	params.Set("fr", "pc")
	params.Set("uc_param_str", "")
	params.Set("pwd_id", pwdID)
	params.Set("stoken", stoken)
	params.Set("pdir_fid", "0")
	params.Set("force", "0")
	params.Set("_page", "1")
	params.Set("_size", "50")
	params.Set("_fetch_banner", "1")
	params.Set("_fetch_share", "1")
	params.Set("_fetch_total", "1")
	params.Set("_sort", "file_type:asc,updated_at:desc")
	params.Set("__dt", "1589")
	params.Set("__t", fmt.Sprintf("%d", now))

	detailURL := "https://drive-h.quark.cn/1/clouddrive/share/sharepage/detail?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", detailURL, nil)
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

	var detailResp QuarkShareDetail
	if err := json.Unmarshal(body, &detailResp); err != nil {
		return nil, err
	}

	if detailResp.Data.List == nil {
		return &response.ShareInfoResponse{
			List:     []response.FileItem{},
			FileSize: 0,
		}, nil
	}

	fileItems := make([]response.FileItem, 0)
	for _, item := range detailResp.Data.List {
		if item.FID != "" {
			fileItems = append(fileItems, response.FileItem{
				FileID:      item.FID,
				FileName:    item.FileName,
				FileType:    item.FileType,
				FileIDToken: item.ShareFidToken,
			})
		}
	}

	return &response.ShareInfoResponse{
		List:     fileItems,
		PwdID:    pwdID,
		Stoken:   stoken,
		FileSize: detailResp.Data.Share.Size,
	}, nil
}

// GetFolderList 获取文件夹列表
func (s *QuarkService) GetFolderList(ctx context.Context, parentID string) (*response.FolderListResponse, error) {
	if parentID == "" {
		parentID = "0"
	}

	now := time.Now().UnixMilli()
	params := url.Values{}
	params.Set("pr", "ucpro")
	params.Set("fr", "pc")
	params.Set("uc_param_str", "")
	params.Set("pdir_fid", parentID)
	params.Set("_page", "1")
	params.Set("_size", "100")
	params.Set("_fetch_total", "false")
	params.Set("_fetch_sub_dirs", "1")
	params.Set("_sort", "")
	params.Set("__dt", "2093126")
	params.Set("__t", fmt.Sprintf("%d", now))

	apiURL := "https://drive-h.quark.cn/1/clouddrive/file/sort?" + params.Encode()

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

	var result QuarkFolderResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Data.List == nil {
		return &response.FolderListResponse{
			Folders: []response.FolderItem{},
		}, nil
	}

	folders := make([]response.FolderItem, 0)
	for _, item := range result.Data.List {
		if item.FID != "" && item.FileType == 0 {
			folders = append(folders, response.FolderItem{
				CID:  item.FID,
				Name: item.FileName,
				Path: []response.FolderItem{},
			})
		}
	}

	return &response.FolderListResponse{
		Folders: folders,
	}, nil
}

// SaveSharedFile 保存分享文件
func (s *QuarkService) SaveSharedFile(ctx context.Context, shareCode, passcode, folderID string, fileIDs []string, fileTokens []string, pwdID, stoken string) error {
	now := time.Now().UnixMilli()
	saveURL := fmt.Sprintf("https://drive-h.quark.cn/1/clouddrive/share/sharepage/save?pr=ucpro&fr=pc&uc_param_str=&__dt=208097&__t=%d", now)

	saveData := map[string]interface{}{
		"fid_list":       fileIDs,
		"fid_token_list": fileTokens,
		"to_pdir_fid":    folderID,
		"pwd_id":         pwdID,
		"stoken":         stoken,
		"pdir_fid":       "0",
		"scene":          "link",
	}

	saveBody, err := json.Marshal(saveData)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", saveURL, bytes.NewReader(saveBody))
	if err != nil {
		return err
	}

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

	var result QuarkSaveResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if result.Status != 200 {
		return fmt.Errorf("保存文件失败: %s", result.Message)
	}

	return nil
}

// setHeaders 设置请求头
func (s *QuarkService) setHeaders(req *http.Request) {
	req.Header.Set("Cookie", s.cookie)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-Ch-Ua", `"Microsoft Edge";v="131", "Chromium";v="131", "Not_A Brand";v="24"`)
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
}
