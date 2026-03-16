package model

import (
	"time"

	"github.com/google/uuid"
)

// DiskInfo 磁盘信息
type DiskInfo struct {
	Name          string  `json:"name"`
	Percent       float64 `json:"percent"`
	TotalGB       float64 `json:"total_gb"`
	AvailableGB   float64 `json:"available_gb"`
}

// ResourceHistory 资源历史模型
type ResourceHistory struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	StoreID           uuid.UUID `gorm:"type:uuid;not null" json:"store_id"`
	CPUPercent        float64   `gorm:"type:decimal(5,2);not null" json:"cpu_percent"`
	MemoryPercent     float64   `gorm:"type:decimal(5,2);not null" json:"memory_percent"`
	MemoryAvailableGB float64   `gorm:"type:decimal(10,2);not null" json:"memory_available_gb"`
	DiskData          []DiskInfo `gorm:"type:jsonb;serializer:json" json:"disk_data"`
	RecordedAt        time.Time `gorm:"not null;default:now()" json:"recorded_at"`

	// 关联
	Store *Store `gorm:"foreignKey:StoreID" json:"store,omitempty"`
}

func (ResourceHistory) TableName() string {
	return "resource_history"
}