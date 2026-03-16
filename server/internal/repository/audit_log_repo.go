package repository

import (
	"store-ops-server/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

func (r *AuditLogRepository) Create(log *model.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *AuditLogRepository) FindAll(storeID *uuid.UUID, actionType string, startDate, endDate *time.Time, page, pageSize int) ([]model.AuditLogWithStore, int64, error) {
	var logs []model.AuditLogWithStore
	var total int64

	query := r.db.Table("audit_logs al").
		Select("al.*, s.name as store_name").
		Joins("LEFT JOIN stores s ON al.store_id = s.id")

	if storeID != nil {
		query = query.Where("al.store_id = ?", storeID)
	}
	if actionType != "" {
		query = query.Where("al.action_type = ?", actionType)
	}
	if startDate != nil {
		query = query.Where("al.created_at >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("al.created_at <= ?", endDate)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("al.created_at DESC").Offset(offset).Limit(pageSize).Scan(&logs).Error
	return logs, total, err
}

func (r *AuditLogRepository) FindByID(id uuid.UUID) (*model.AuditLog, error) {
	var log model.AuditLog
	err := r.db.Where("id = ?", id).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *AuditLogRepository) DeleteOlderThan(days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)
	return r.db.Where("created_at < ?", cutoff).Delete(&model.AuditLog{}).Error
}