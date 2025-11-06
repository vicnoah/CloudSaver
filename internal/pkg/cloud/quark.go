package cloud

import (
	"context"
	"cloudsaver/internal/api/response"
	"cloudsaver/internal/repository"
	"errors"
)

// QuarkService 夸克网盘服务
type QuarkService struct {
	settingRepo *repository.SettingRepository
	cookie      string
}

// NewQuarkService 创建夸克网盘服务
func NewQuarkService(settingRepo *repository.SettingRepository) *QuarkService {
	return &QuarkService{
		settingRepo: settingRepo,
	}
}

// SetCookie 设置Cookie
func (s *QuarkService) SetCookie(userUUID string) error {
	setting, err := s.settingRepo.GetUserSetting(userUUID)
	if err != nil {
		return err
	}

	if setting.QuarkCookie == "" {
		return errors.New("请先设置夸克网盘Cookie")
	}

	s.cookie = setting.QuarkCookie
	return nil
}

// GetShareInfo 获取分享信息
func (s *QuarkService) GetShareInfo(ctx context.Context, shareCode string, passcode string) (*response.ShareInfoResponse, error) {
	// TODO: 实现夸克网盘API调用
	// 步骤1: POST https://drive-h.quark.cn/1/clouddrive/share/sharepage/token
	// 步骤2: GET https://drive-h.quark.cn/1/clouddrive/share/sharepage/detail
	return &response.ShareInfoResponse{
		List:     []response.FileItem{},
		FileSize: 0,
	}, nil
}

// GetFolderList 获取文件夹列表
func (s *QuarkService) GetFolderList(ctx context.Context, parentID string) (*response.FolderListResponse, error) {
	// TODO: 实现夸克网盘API调用
	// GET https://drive-h.quark.cn/1/clouddrive/file/sort
	return &response.FolderListResponse{
		Folders: []response.FolderItem{},
	}, nil
}

// SaveSharedFile 保存分享文件
func (s *QuarkService) SaveSharedFile(ctx context.Context, shareCode, passcode, folderID string, fileIDs []string, fileTokens []string, pwdID, stoken string) error {
	// TODO: 实现夸克网盘API调用
	// POST https://drive-h.quark.cn/1/clouddrive/share/sharepage/save
	return nil
}
