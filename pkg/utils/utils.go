package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/google/uuid"
)

// GenerateUUID 生成新的 UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateShortID 生成短 ID
func GenerateShortID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// ValidateIP 验证 IP 地址
func ValidateIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// ValidatePort 验证端口号
func ValidatePort(port int) bool {
	return port > 0 && port <= 65535
}

// FormatDuration 格式化时间间隔
func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d)/float64(time.Millisecond))
	}
	if d < time.Minute {
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.2fm", d.Minutes())
	}
	return fmt.Sprintf("%.2fh", d.Hours())
}

// ParseHostPort 解析主机和端口
func ParseHostPort(addr string) (host string, port string, err error) {
	host, port, err = net.SplitHostPort(addr)
	if err != nil {
		// 如果没有端口，尝试添加默认端口
		if strings.Contains(err.Error(), "missing port") {
			return addr, "80", nil
		}
		return "", "", err
	}
	return host, port, nil
}

// IsValidUUID 验证 UUID 格式
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// RandomBytes 生成随机字节
func RandomBytes(n int) []byte {
	bytes := make([]byte, n)
	rand.Read(bytes)
	return bytes
}

// BytesToHex 字节转十六进制字符串
func BytesToHex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

// HexToBytes 十六进制字符串转字节
func HexToBytes(s string) ([]byte, error) {
	return hex.DecodeString(s)
}

// Contains 检查字符串切片是否包含指定字符串
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveDuplicates 移除字符串切片中的重复项
func RemoveDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var result []string
	
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	
	return result
}
