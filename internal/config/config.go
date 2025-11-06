package config

import (
	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Telegram TelegramConfig `mapstructure:"telegram"`
	Proxy    ProxyConfig    `mapstructure:"proxy"`
	Log      LogConfig      `mapstructure:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Path     string `mapstructure:"path"`
	LogLevel string `mapstructure:"log_level"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

// TelegramConfig Telegram配置
type TelegramConfig struct {
	BaseURL  string            `mapstructure:"base_url"`
	Channels []TelegramChannel `mapstructure:"channels"`
}

// TelegramChannel Telegram频道
type TelegramChannel struct {
	ID   string `mapstructure:"id"`
	Name string `mapstructure:"name"`
}

// ProxyConfig 代理配置
type ProxyConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	File       string `mapstructure:"file"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

// CloudPatterns 云盘链接正则模式
type CloudPatterns struct {
	Pan115 string `mapstructure:"pan115"`
	Quark  string `mapstructure:"quark"`
	Baidu  string `mapstructure:"baidu"`
	Aliyun string `mapstructure:"aliyun"`
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// 设置默认值
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "release")
	viper.SetDefault("database.path", "./data/cloudsaver.db")
	viper.SetDefault("database.log_level", "warn")
	viper.SetDefault("jwt.secret", "your-secret-key-change-in-production")
	viper.SetDefault("jwt.expire_hours", 168)

	// 环境变量覆盖
	viper.AutomaticEnv()
	viper.SetEnvPrefix("CLOUDSAVER")
	// 绑定环境变量到配置项
	viper.BindEnv("server.port", "CLOUDSAVER_SERVER_PORT")
	viper.BindEnv("server.mode", "CLOUDSAVER_SERVER_MODE")
	viper.BindEnv("database.path", "CLOUDSAVER_DATABASE_PATH")

	if err := viper.ReadInConfig(); err != nil {
		// 配置文件不存在时使用默认值
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
