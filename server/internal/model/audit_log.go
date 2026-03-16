package model

import (
	"time"

	"github.com/google/uuid"
)

// ActionType 操作类型
type ActionType string

const (
	ActionTypeDesktopConnect    ActionType = "desktop_connect"
	ActionTypeDesktopDisconnect ActionType = "desktop_disconnect"
	ActionTypeCommandExecute    ActionType = "command_execute"
	ActionTypeSoftwarePush      ActionType = "software_push"
	ActionTypeProcessKill       ActionType = "process_kill"
	ActionTypeServiceControl    ActionType = "service_control"
)

// AuditLog 审计日志模型
type AuditLog struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	StoreID    *uuid.UUID `gorm:"type:uuid" json:"store_id"`
	ActionType string     `gorm:"type:varchar(50);not null" json:"action_type"`
	Action     string     `gorm:"type:varchar(100);not null" json:"action"`
	Detail     string     `gorm:"type:text" json:"detail"`
	Operator   string     `gorm:"type:varchar(50);not null" json:"operator"`
	CreatedAt  time.Time  `gorm:"not null;default:now()" json:"created_at"`

	// 关联
	Store *Store `gorm:"foreignKey:StoreID" json:"store,omitempty"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

// AuditLogWithStore 带门店信息的审计日志
type AuditLogWithStore struct {
	ID         uuid.UUID  `json:"id"`
	StoreID    *uuid.UUID `json:"store_id"`
	ActionType string     `json:"action_type"`
	Action     string     `json:"action"`
	Detail     string     `json:"detail"`
	Operator   string     `json:"operator"`
	CreatedAt  time.Time  `json:"created_at"`
	StoreName  string     `json:"store_name,omitempty"`
}