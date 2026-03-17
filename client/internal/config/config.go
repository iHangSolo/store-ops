package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	ServerURL         string `json:"server_url"`
	Token             string `json:"token"`
	StoreName         string `json:"store_name"`
	DeviceID          string `json:"device_id"`
	DeviceFingerprint string `json:"device_fingerprint"`
	RustDeskID        string `json:"rustdesk_id"`
}

var (
	cfg     *Config
	cfgOnce sync.Once
	cfgPath string
)

func Init() error {
	var initErr error
	cfgOnce.Do(func() {
		// 获取配置文件路径
		exePath, err := os.Executable()
		if err != nil {
			initErr = err
			return
		}
		cfgPath = filepath.Join(filepath.Dir(exePath), "config.json")

		// 加载配置
		cfg = &Config{
			ServerURL: "wss://yhzwj.top/ws",
		}

		data, err := os.ReadFile(cfgPath)
		if err == nil {
			json.Unmarshal(data, cfg)
		}
	})
	return initErr
}

func Get() *Config {
	return cfg
}

func Save() error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cfgPath, data, 0644)
}

func SetStoreName(name string) error {
	cfg.StoreName = name
	return Save()
}

func SetToken(token string) error {
	cfg.Token = token
	return Save()
}

func SetDeviceID(id string) error {
	cfg.DeviceID = id
	return Save()
}

func SetDeviceFingerprint(fp string) error {
	cfg.DeviceFingerprint = fp
	return Save()
}

func SetRustDeskID(id string) error {
	cfg.RustDeskID = id
	return Save()
}