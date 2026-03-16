package model

import (
	"time"

	"github.com/google/uuid"
)

// DesktopSessionStatus 桌面会话状态
type DesktopSessionStatus string

const (
	DesktopSessionActive    DesktopSessionStatus = "active"
	DesktopSessionCompleted DesktopSessionStatus = "completed"
)

// DesktopSession 远程桌面会话模型
type DesktopSession struct {
	ID              uuid.UUID            `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	StoreID         uuid.UUID            `gorm:"type:uuid;not null" json:"store_id"`
	Operator        string               `gorm:"type:varchar(50);not null" json:"operator"`
	StartedAt       time.Time            `gorm:"not null" json:"started_at"`
	EndedAt         *time.Time           `json:"ended_at"`
	DurationSeconds int                  `json:"duration_seconds"`
	Status          DesktopSessionStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"`

	// 关联
	Store *Store `gorm:"foreignKey:StoreID" json:"store,omitempty"`
}

func (DesktopSession) TableName() string {
	return "desktop_sessions"
}

// DesktopQueue 远程桌面排队模型
type DesktopQueue struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	StoreID   uuid.UUID `gorm:"type:uuid;not null" json:"store_id"`
	Operator  string    `gorm:"type:varchar(50);not null" json:"operator"`
	Position  int       `gorm:"not null" json:"position"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`

	// 关联
	Store *Store `gorm:"foreignKey:StoreID" json:"store,omitempty"`
}

func (DesktopQueue) TableName() string {
	return "desktop_queue"
}