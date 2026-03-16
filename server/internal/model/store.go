package model

import (
	"time"

	"github.com/google/uuid"
)

// StoreStatus 门店状态
type StoreStatus string

const (
	StoreStatusPending  StoreStatus = "pending"
	StoreStatusApproved StoreStatus = "approved"
	StoreStatusRejected StoreStatus = "rejected"
)

// Store 门店模型
type Store struct {
	ID                uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name              string      `gorm:"type:varchar(100);not null" json:"name"`
	DeviceID          uuid.UUID   `gorm:"type:uuid;uniqueIndex;not null" json:"device_id"`
	DeviceFingerprint string      `gorm:"type:varchar(64)" json:"device_fingerprint"`
	Token             string      `gorm:"type:varchar(64);uniqueIndex;not null" json:"token"`
	RustDeskID        string      `gorm:"type:varchar(20)" json:"rustdesk_id"`
	ClientVersion     string      `gorm:"type:varchar(20)" json:"client_version"`
	Status            StoreStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	IsOnline          bool        `gorm:"not null;default:false" json:"is_online"`
	CPUPercent        float64     `gorm:"type:decimal(5,2)" json:"cpu_percent"`
	MemoryPercent     float64     `gorm:"type:decimal(5,2)" json:"memory_percent"`
	LastSeen          *time.Time  `json:"last_seen"`
	CreatedAt         time.Time   `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt         time.Time   `gorm:"not null;default:now()" json:"updated_at"`
}

func (Store) TableName() string {
	return "stores"
}

// StoreWithResource 带资源信息的门店
type StoreWithResource struct {
	Store
	Resource *ResourceSummary `json:"resource,omitempty"`
}

// ResourceSummary 资源摘要
type ResourceSummary struct {
	CPUPercent      float64 `json:"cpu_percent"`
	MemoryPercent   float64 `json:"memory_percent"`
	MemoryAvailable float64 `json:"memory_available_gb"`
	DiskPercent     float64 `json:"disk_percent"`
	DiskAvailable   float64 `json:"disk_available_gb"`
}