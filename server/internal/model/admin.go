package model

import (
	"time"

	"github.com/google/uuid"
)

// AdminStatus 管理员状态
type AdminStatus string

const (
	AdminStatusActive   AdminStatus = "active"
	AdminStatusInactive AdminStatus = "inactive"
)

// Admin 管理员模型
type Admin struct {
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Username     string      `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password     string      `gorm:"type:varchar(255);not null" json:"-"`
	Nickname     string      `gorm:"type:varchar(50)" json:"nickname"`
	Role         string      `gorm:"type:varchar(20);not null;default:'admin'" json:"role"`
	Status       AdminStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	LastLogin    *time.Time  `json:"last_login"`
	CreatedAt    time.Time   `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time   `gorm:"not null;default:now()" json:"updated_at"`
}

func (Admin) TableName() string {
	return "admins"
}