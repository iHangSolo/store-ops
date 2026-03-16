package model

import (
	"time"

	"github.com/google/uuid"
)

// CommandStatus 命令状态
type CommandStatus string

const (
	CommandStatusExecuting CommandStatus = "executing"
	CommandStatusSuccess   CommandStatus = "success"
	CommandStatusFailed    CommandStatus = "failed"
	CommandStatusTimeout   CommandStatus = "timeout"
)

// CommandLog 命令日志模型
type CommandLog struct {
	ID         uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	StoreID    *uuid.UUID    `gorm:"type:uuid" json:"store_id"`
	Operator   string        `gorm:"type:varchar(50);not null" json:"operator"`
	Command    string        `gorm:"type:text;not null" json:"command"`
	Output     string        `gorm:"type:text" json:"output"`
	Status     CommandStatus `gorm:"type:varchar(20);not null" json:"status"`
	ExecutedAt time.Time     `gorm:"not null;default:now()" json:"executed_at"`
	DurationMs int           `json:"duration_ms"`

	// 关联
	Store *Store `gorm:"foreignKey:StoreID" json:"store,omitempty"`
}

func (CommandLog) TableName() string {
	return "command_logs"
}