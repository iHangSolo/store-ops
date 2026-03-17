package main

import (
	"context"
	"log"

	"store-ops-client/internal/config"
	"store-ops-client/internal/connection"
	"store-ops-client/internal/device"
)

// App 主应用结构
type App struct {
	ctx     context.Context
	version string
}

// NewApp 创建新的 App 实例
func NewApp() *App {
	return &App{
		version: "1.0.0",
	}
}

// startup 应用启动时调用
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	log.Println("门店运维客户端启动...")

	// 初始化配置
	if err := config.Init(); err != nil {
		log.Printf("初始化配置失败: %v", err)
		return
	}

	cfg := config.Get()

	// 如果没有设备ID，生成新的
	if cfg.DeviceID == "" {
		deviceID := device.GenerateDeviceID()
		config.SetDeviceID(deviceID)
	}

	// 如果没有设备指纹，生成新的
	if cfg.DeviceFingerprint == "" {
		fp := device.GenerateFingerprint()
		config.SetDeviceFingerprint(fp)
	}

	// 初始化连接
	if err := connection.Init(); err != nil {
		log.Printf("初始化连接失败: %v", err)
		return
	}

	// 如果有TOKEN，尝试连接
	if cfg.Token != "" {
		client := connection.GetClient()
		if err := client.Connect(); err != nil {
			log.Printf("连接服务器失败: %v", err)
		}
	}
}

// shutdown 应用关闭时调用
func (a *App) shutdown(ctx context.Context) {
	log.Println("门店运维客户端关闭...")
	client := connection.GetClient()
	if client != nil {
		client.Disconnect()
	}
}

// GetVersion 获取客户端版本
func (a *App) GetVersion() string {
	return a.version
}

// GetConnectionStatus 获取连接状态
func (a *App) GetConnectionStatus() string {
	client := connection.GetClient()
	if client == nil {
		return "disconnected"
	}
	if client.IsConnected() {
		return "connected"
	}
	return "disconnected"
}

// GetStoreName 获取门店名称
func (a *App) GetStoreName() string {
	return config.Get().StoreName
}

// SetStoreName 设置门店名称
func (a *App) SetStoreName(name string) error {
	return config.SetStoreName(name)
}

// GetToken 获取TOKEN
func (a *App) GetToken() string {
	return config.Get().Token
}

// SetToken 设置TOKEN
func (a *App) SetToken(token string) error {
	err := config.SetToken(token)
	if err != nil {
		return err
	}

	// 重新连接
	client := connection.GetClient()
	if client != nil {
		client.Disconnect()
		if err := client.Connect(); err != nil {
			log.Printf("连接服务器失败: %v", err)
		}
	}
	return nil
}

// GetDeviceID 获取设备ID
func (a *App) GetDeviceID() string {
	return config.Get().DeviceID
}

// GetRustDeskID 获取RustDesk ID
func (a *App) GetRustDeskID() string {
	return config.Get().RustDeskID
}