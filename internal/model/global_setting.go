package model

import "time"

// GlobalSetting 全局设置模型
type GlobalSetting struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	HTTPProxyHost   string    `gorm:"type:varchar(255);not null;default:'127.0.0.1'" json:"http_proxy_host"`
	HTTPProxyPort   int       `gorm:"type:int;not null;default:7890" json:"http_proxy_port"`
	IsProxyEnabled  bool      `gorm:"type:boolean;not null;default:true" json:"is_proxy_enabled"`
	CommonUserCode  int       `gorm:"type:int;default:9527" json:"common_user_code"`
	AdminUserCode   int       `gorm:"type:int;not null;default:230713" json:"admin_user_code"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TableName 指定表名
func (GlobalSetting) TableName() string {
	return "global_settings"
}
