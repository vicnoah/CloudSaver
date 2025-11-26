package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID      string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"uuid"`
	Username  string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Role      int            `gorm:"type:int;not null;default:0" json:"role"` // 0-普通用户 1-管理员
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
