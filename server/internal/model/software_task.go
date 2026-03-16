package model

import (
	"time"

	"github.com/google/uuid"
)

// SoftwareStatus 软件任务状态
type SoftwareStatus string

const (
	SoftwareStatusPending    SoftwareStatus = "pending"
	SoftwareStatusUploading  SoftwareStatus = "uploading"
	SoftwareStatusDownloading SoftwareStatus = "downloading"
	SoftwareStatusInstalling SoftwareStatus = "installing"
	SoftwareStatusInstalled  SoftwareStatus = "installed"
	SoftwareStatusFailed     SoftwareStatus = "failed"
)

// SoftwareTask 软件任务模型
type SoftwareTask struct {
	ID          uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	StoreID     uuid.UUID     `gorm:"type:uuid;not null" json:"store_id"`
	Operator    string        `gorm:"type:varchar(50);not null" json:"operator"`
	FileName    string        `gorm:"type:varchar(255);not null" json:"file_name"`
	FileSize    int64         `json:"file_size"`
	InstallArgs string        `gorm:"type:varchar(500)" json:"install_args"`
	Status      SoftwareStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Message     string        `json:"message"`
	CreatedAt   time.Time     `gorm:"not null;default:now()" json:"created_at"`
	CompletedAt *time.Time    `json:"completed_at"`

	// 关联
	Store *Store `gorm:"foreignKey:StoreID" json:"store,omitempty"`
}

func (SoftwareTask) TableName() string {
	return "software_tasks"
}