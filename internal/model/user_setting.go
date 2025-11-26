package model

import "time"

// UserSetting 用户设置模型
type UserSetting struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserUUID        string    `gorm:"type:varchar(36);uniqueIndex;not null" json:"user_uuid"`
	Cloud115Cookie  string    `gorm:"type:text" json:"cloud115_cookie"`
	QuarkCookie     string    `gorm:"type:text" json:"quark_cookie"`
	Cloud115UserID  string    `gorm:"type:varchar(100)" json:"cloud115_user_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TableName 指定表名
func (UserSetting) TableName() string {
	return "user_settings"
}
