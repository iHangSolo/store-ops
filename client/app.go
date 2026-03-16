package main

import (
	"context"
	"fmt"
	"log"
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

	// TODO: 初始化配置
	// TODO: 初始化 SQLite
	// TODO: 连接总部服务器
}

// shutdown 应用关闭时调用
func (a *App) shutdown(ctx context.Context) {
	log.Println("门店运维客户端关闭...")
	// TODO: 断开连接
	// TODO: 清理资源
}

// GetVersion 获取客户端版本
func (a *App) GetVersion() string {
	return a.version
}

// GetConnectionStatus 获取连接状态
func (a *App) GetConnectionStatus() string {
	// TODO: 返回实际连接状态
	return "disconnected"
}

// GetStoreName 获取门店名称
func (a *App) GetStoreName() string {
	// TODO: 从本地存储读取
	return ""
}

// SetStoreName 设置门店名称
func (a *App) SetStoreName(name string) error {
	// TODO: 保存到本地存储
	fmt.Println("门店名称设置为:", name)
	return nil
}