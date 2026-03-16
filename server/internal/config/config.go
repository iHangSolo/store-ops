package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	RustDesk RustDeskConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type RustDeskConfig struct {
	Server string
}

var cfg *Config

// Load 加载配置
func Load() (*Config, error) {
	// 加载 .env 文件
	godotenv.Load()

	cfg = &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("SERVER_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "storeops"),
			Password: getEnv("DB_PASSWORD", "storeops123"),
			DBName:   getEnv("DB_NAME", "store_ops"),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "your_jwt_secret"),
			ExpireHours: getEnvInt("JWT_EXPIRE_HOURS", 24),
		},
		RustDesk: RustDeskConfig{
			Server: getEnv("RUSTDESK_SERVER", "localhost:21115"),
		},
	}

	return cfg, nil
}

// Get 获取配置
func Get() *Config {
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}