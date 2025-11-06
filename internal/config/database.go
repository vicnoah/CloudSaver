package config

import (
	"fmt"

	"cloudsaver/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDatabase 初始化数据库
func InitDatabase(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(
		&model.User{},
		&model.GlobalSetting{},
		&model.UserSetting{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// 初始化默认全局设置
	var count int64
	db.Model(&model.GlobalSetting{}).Count(&count)
	if count == 0 {
		defaultSetting := &model.GlobalSetting{
			HTTPProxyHost:  "127.0.0.1",
			HTTPProxyPort:  7890,
			IsProxyEnabled: false,
			CommonUserCode: 9527,
			AdminUserCode:  230713,
		}
		if err := db.Create(defaultSetting).Error; err != nil {
			return nil, fmt.Errorf("failed to create default settings: %w", err)
		}
	}

	return db, nil
}
