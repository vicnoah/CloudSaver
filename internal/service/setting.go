package service

import (
	"cloudsaver/internal/model"
	"cloudsaver/internal/repository"
)

// SettingService 设置服务
type SettingService struct {
	settingRepo *repository.SettingRepository
}

// NewSettingService 创建设置服务
func NewSettingService(settingRepo *repository.SettingRepository) *SettingService {
	return &SettingService{
		settingRepo: settingRepo,
	}
}

// GetUserSetting 获取用户设置
func (s *SettingService) GetUserSetting(userUUID string) (*model.UserSetting, error) {
	return s.settingRepo.GetUserSetting(userUUID)
}

// SaveUserSetting 保存用户设置
func (s *SettingService) SaveUserSetting(userUUID string, cloud115Cookie, quarkCookie string) error {
	setting, err := s.settingRepo.GetUserSetting(userUUID)
	if err != nil {
		return err
	}

	if cloud115Cookie != "" {
		setting.Cloud115Cookie = cloud115Cookie
	}

	if quarkCookie != "" {
		setting.QuarkCookie = quarkCookie
	}

	return s.settingRepo.SaveUserSetting(setting)
}
