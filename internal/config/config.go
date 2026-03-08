package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Wechat WechatConfig `mapstructure:"wechat"`
	Server ServerConfig `mapstructure:"server"`
	API    APIConfig    `mapstructure:"api"`
	Log    LogConfig    `mapstructure:"log"`
}

type WechatConfig struct {
	AppID     string `mapstructure:"app_id"`
	AppSecret string `mapstructure:"app_secret"`
	Token     string `mapstructure:"token"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// API 访问控制配置
type APIConfig struct {
	IPWhitelist []string `mapstructure:"ip_whitelist"`
	APIKeys     []string `mapstructure:"api_keys"`
}

// 日志配置
type LogConfig struct {
	Level       string `mapstructure:"level"`       // 日志级别：debug, info, warn, error
	Path        string `mapstructure:"path"`        // 日志文件路径
	Console     bool   `mapstructure:"console"`     // 是否在控制台打印日志
	MaxSize     int    `mapstructure:"max_size"`    // 单个日志文件最大大小（MB）
	MaxBackups  int    `mapstructure:"max_backups"` // 最大备份文件数
	MaxAge      int    `mapstructure:"max_age"`     // 日志文件最大保存天数
	Compress    bool   `mapstructure:"compress"`    // 是否压缩日志文件
}

func LoadConfig() (*Config, error) {
	// 设置默认配置文件路径
	configDir := "./configs"
	configFile := "config.yaml"

	// 检查配置文件是否存在
	configPath := filepath.Join(configDir, configFile)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %s", configPath)
	}

	// 读取配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 绑定环境变量
	viper.AutomaticEnv()

	// 优先使用环境变量
	if appID := os.Getenv("WECHAT_APPID"); appID != "" {
		viper.Set("wechat.app_id", appID)
	}
	if appSecret := os.Getenv("WECHAT_APPSECRET"); appSecret != "" {
		viper.Set("wechat.app_secret", appSecret)
	}
	if token := os.Getenv("WECHAT_TOKEN"); token != "" {
		viper.Set("wechat.token", token)
	}

	// 解析配置
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	return &config, nil
}

