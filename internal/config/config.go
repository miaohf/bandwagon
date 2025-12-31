package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config 主配置结构
type Config struct {
	Port     int           `json:"port"`
	LogLevel string        `json:"log_level"`
	Inbounds []InboundConfig `json:"inbounds"`
}

// InboundConfig 入站配置
type InboundConfig struct {
	Protocol string      `json:"protocol"`
	Port     int         `json:"port"`
	Settings interface{} `json:"settings"`
	StreamSettings *StreamSettings `json:"stream_settings,omitempty"`
}

// StreamSettings 流设置
type StreamSettings struct {
	Network  string             `json:"network"`
	Security string             `json:"security"`
	TLSSettings *TLSSettings   `json:"tls_settings,omitempty"`
	RealitySettings *RealitySettings `json:"reality_settings,omitempty"`
}

// TLSSettings TLS 配置
type TLSSettings struct {
	ServerName string   `json:"server_name"`
	Certificates []Certificate `json:"certificates"`
}

// Certificate 证书配置
type Certificate struct {
	CertificateFile string `json:"certificate_file"`
	KeyFile         string `json:"key_file"`
}

// RealitySettings Reality 配置
type RealitySettings struct {
	Show         bool     `json:"show"`
	Dest         string   `json:"dest"`
	XVer         int      `json:"xver"`
	ServerNames  []string `json:"server_names"`
	PrivateKey   string   `json:"private_key"`
	MinClientVer string   `json:"min_client_ver"`
	MaxClientVer string   `json:"max_client_ver"`
	MaxTimeDiff  int      `json:"max_time_diff"`
	ShortIds     []string `json:"short_ids"`
}

// VLESSSettings VLESS 协议设置
type VLESSSettings struct {
	Clients []VLESSClient `json:"clients"`
	Decryption string     `json:"decryption"`
	Fallbacks  []Fallback `json:"fallbacks,omitempty"`
}

// VLESSClient VLESS 客户端配置
type VLESSClient struct {
	ID    string `json:"id"`
	Flow  string `json:"flow,omitempty"`
	Email string `json:"email,omitempty"`
}

// Fallback 回落配置
type Fallback struct {
	Name string `json:"name,omitempty"`
	Alpn string `json:"alpn,omitempty"`
	Path string `json:"path,omitempty"`
	Dest string `json:"dest"`
	Xver int    `json:"xver,omitempty"`
}

// LoadConfig 加载配置文件
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("配置验证失败: %v", err)
	}

	return &config, nil
}

// validateConfig 验证配置
func validateConfig(config *Config) error {
	if config.Port <= 0 || config.Port > 65535 {
		return fmt.Errorf("无效的端口号: %d", config.Port)
	}

	if len(config.Inbounds) == 0 {
		return fmt.Errorf("至少需要一个入站配置")
	}

	for i, inbound := range config.Inbounds {
		if inbound.Protocol == "" {
			return fmt.Errorf("入站配置 %d 缺少协议", i)
		}
		if inbound.Port <= 0 || inbound.Port > 65535 {
			return fmt.Errorf("入站配置 %d 端口号无效: %d", i, inbound.Port)
		}
	}

	return nil
}

// SaveConfig 保存配置文件
func SaveConfig(config *Config, filename string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
}
