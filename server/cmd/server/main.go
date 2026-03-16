package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// 获取配置
	port := getEnv("SERVER_PORT", "8080")

	log.Printf("Starting store-ops server on port %s...", port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}