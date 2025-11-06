package config

import (
	"fmt"
	"log"

	"cloudsaver/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDatabase 初始化数据库并执行自动迁移
func InitDatabase(dbPath string) (*gorm.DB, error) {
	// 配置 GORM
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	log.Println("开始数据库迁移...")

	// 使用 GORM AutoMigrate 自动迁移所有表结构
	// AutoMigrate 会自动创建表、缺失的列和索引，但不会删除未使用的列
	if err := db.AutoMigrate(
		&model.User{},
		&model.GlobalSetting{},
		&model.UserSetting{},
	); err != nil {
		return nil, fmt.Errorf("failed to auto migrate database: %w", err)
	}

	log.Println("数据库迁移完成")

	// 初始化默认数据
	if err := initDefaultData(db); err != nil {
		return nil, fmt.Errorf("failed to initialize default data: %w", err)
	}

	return db, nil
}

// initDefaultData 初始化默认数据
func initDefaultData(db *gorm.DB) error {
	// 检查并创建默认全局设置
	var count int64
	db.Model(&model.GlobalSetting{}).Count(&count)
	if count == 0 {
		log.Println("初始化默认全局设置...")
		defaultSetting := &model.GlobalSetting{
			HTTPProxyHost:  "127.0.0.1",
			HTTPProxyPort:  7890,
			IsProxyEnabled: false,
			CommonUserCode: 9527,
			AdminUserCode:  230713,
		}
		if err := db.Create(defaultSetting).Error; err != nil {
			return fmt.Errorf("failed to create default global settings: %w", err)
		}
		log.Println("默认全局设置创建成功")
	}

	return nil
}
