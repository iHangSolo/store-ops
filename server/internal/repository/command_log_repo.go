package repository

import (
	"store-ops-server/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommandLogRepository struct {
	db *gorm.DB
}

func NewCommandLogRepository(db *gorm.DB) *CommandLogRepository {
	return &CommandLogRepository{db: db}
}

func (r *CommandLogRepository) Create(log *model.CommandLog) error {
	return r.db.Create(log).Error
}

func (r *CommandLogRepository) Update(log *model.CommandLog) error {
	return r.db.Save(log).Error
}

func (r *CommandLogRepository) FindByID(id uuid.UUID) (*model.CommandLog, error) {
	var log model.CommandLog
	err := r.db.Where("id = ?", id).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *CommandLogRepository) FindByStoreID(storeID uuid.UUID, page, pageSize int) ([]model.CommandLog, int64, error) {
	var logs []model.CommandLog
	var total int64

	query := r.db.Model(&model.CommandLog{}).Where("store_id = ?", storeID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("executed_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error
	return logs, total, err
}

func (r *CommandLogRepository) FindByOperator(operator string, page, pageSize int) ([]model.CommandLog, int64, error) {
	var logs []model.CommandLog
	var total int64

	query := r.db.Model(&model.CommandLog{}).Where("operator = ?", operator)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("executed_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error
	return logs, total, err
}