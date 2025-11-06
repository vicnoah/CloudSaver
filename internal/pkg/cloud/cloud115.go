package cloud

import (
	"context"
	"cloudsaver/internal/api/response"
	"cloudsaver/internal/repository"
	"errors"
)

// Cloud115Service 115网盘服务
type Cloud115Service struct {
	settingRepo *repository.SettingRepository
	cookie      string
}

// NewCloud115Service 创建115网盘服务
func NewCloud115Service(settingRepo *repository.SettingRepository) *Cloud115Service {
	return &Cloud115Service{
		settingRepo: settingRepo,
	}
}

// SetCookie 设置Cookie
func (s *Cloud115Service) SetCookie(userUUID string) error {
	setting, err := s.settingRepo.GetUserSetting(userUUID)
	if err != nil {
		return err
	}

	if setting.Cloud115Cookie == "" {
		return errors.New("请先设置115网盘Cookie")
	}

	s.cookie = setting.Cloud115Cookie
	return nil
}

// GetShareInfo 获取分享信息
func (s *Cloud115Service) GetShareInfo(ctx context.Context, shareCode string, passcode string) (*response.ShareInfoResponse, error) {
	// TODO: 实现115网盘API调用
	// 这里需要调用 https://webapi.115.com/share/snap
	return &response.ShareInfoResponse{
		List:     []response.FileItem{},
		FileSize: 0,
	}, nil
}

// GetFolderList 获取文件夹列表
func (s *Cloud115Service) GetFolderList(ctx context.Context, parentID string) (*response.FolderListResponse, error) {
	// TODO: 实现115网盘API调用
	// 这里需要调用 https://webapi.115.com/files
	return &response.FolderListResponse{
		Folders: []response.FolderItem{},
	}, nil
}

// SaveSharedFile 保存分享文件
func (s *Cloud115Service) SaveSharedFile(ctx context.Context, shareCode, passcode, folderID string, fileIDs []string, fileTokens []string, pwdID, stoken string) error {
	// TODO: 实现115网盘API调用
	// 这里需要调用 https://webapi.115.com/share/receive
	return nil
}
