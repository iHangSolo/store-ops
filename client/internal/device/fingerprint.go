package device

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"runtime"
	"strings"

	"github.com/google/uuid"
)

// GenerateDeviceID 生成设备ID
func GenerateDeviceID() string {
	return uuid.New().String()
}

// GenerateFingerprint 生成设备指纹
func GenerateFingerprint() string {
	var parts []string

	// 主机名
	if hostname, err := os.Hostname(); err == nil {
		parts = append(parts, hostname)
	}

	// 操作系统信息
	parts = append(parts, runtime.GOOS)
	parts = append(parts, runtime.GOARCH)

	// CPU 核心数
	parts = append(parts, string(rune(runtime.NumCPU())))

	// 环境变量
	if user := os.Getenv("USERNAME"); user != "" {
		parts = append(parts, user)
	}
	if user := os.Getenv("USER"); user != "" {
		parts = append(parts, user)
	}

	// 组合并哈希
	combined := strings.Join(parts, "|")
	hash := md5.Sum([]byte(combined))
	return hex.EncodeToString(hash[:])
}