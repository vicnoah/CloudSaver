package repository

import (
	"cloudsaver/internal/model"
	"gorm.io/gorm"
)

// SettingRepository 设置仓库
type SettingRepository struct {
	db *gorm.DB
}

// NewSettingRepository 创建设置仓库
func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

// GetGlobalSetting 获取全局设置
func (r *SettingRepository) GetGlobalSetting() (*model.GlobalSetting, error) {
	var setting model.GlobalSetting
	if err := r.db.First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

// UpdateGlobalSetting 更新全局设置
func (r *SettingRepository) UpdateGlobalSetting(setting *model.GlobalSetting) error {
	return r.db.Save(setting).Error
}

// GetUserSetting 获取用户设置
func (r *SettingRepository) GetUserSetting(userUUID string) (*model.UserSetting, error) {
	var setting model.UserSetting
	if err := r.db.Where("user_uuid = ?", userUUID).First(&setting).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果不存在，创建默认设置
			setting = model.UserSetting{
				UserUUID: userUUID,
			}
			if err := r.db.Create(&setting).Error; err != nil {
				return nil, err
			}
			return &setting, nil
		}
		return nil, err
	}
	return &setting, nil
}

// SaveUserSetting 保存用户设置
func (r *SettingRepository) SaveUserSetting(setting *model.UserSetting) error {
	return r.db.Save(setting).Error
}
