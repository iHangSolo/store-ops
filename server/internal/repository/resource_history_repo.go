package repository

import (
	"store-ops-server/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResourceHistoryRepository struct {
	db *gorm.DB
}

func NewResourceHistoryRepository(db *gorm.DB) *ResourceHistoryRepository {
	return &ResourceHistoryRepository{db: db}
}

func (r *ResourceHistoryRepository) Create(history *model.ResourceHistory) error {
	return r.db.Create(history).Error
}

func (r *ResourceHistoryRepository) FindByStoreID(storeID uuid.UUID, startTime, endTime time.Time, page, pageSize int) ([]model.ResourceHistory, int64, error) {
	var history []model.ResourceHistory
	var total int64

	query := r.db.Model(&model.ResourceHistory{}).Where("store_id = ?", storeID)
	if !startTime.IsZero() {
		query = query.Where("recorded_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("recorded_at <= ?", endTime)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("recorded_at DESC").Offset(offset).Limit(pageSize).Find(&history).Error
	return history, total, err
}

func (r *ResourceHistoryRepository) FindLatestByStoreID(storeID uuid.UUID) (*model.ResourceHistory, error) {
	var history model.ResourceHistory
	err := r.db.Where("store_id = ?", storeID).Order("recorded_at DESC").First(&history).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

func (r *ResourceHistoryRepository) DeleteOlderThan(storeID uuid.UUID, days int) error {
	cutoff := time.Now().AddDate(0, 0, -days)
	return r.db.Where("store_id = ? AND recorded_at < ?", storeID, cutoff).Delete(&model.ResourceHistory{}).Error
}

func (r *ResourceHistoryRepository) GetAverageCPU(storeID uuid.UUID, startTime, endTime time.Time) (float64, error) {
	var avg float64
	err := r.db.Model(&model.ResourceHistory{}).
		Where("store_id = ? AND recorded_at >= ? AND recorded_at <= ?", storeID, startTime, endTime).
		Select("AVG(cpu_percent)").
		Scan(&avg).Error
	return avg, err
}

func (r *ResourceHistoryRepository) GetAverageMemory(storeID uuid.UUID, startTime, endTime time.Time) (float64, error) {
	var avg float64
	err := r.db.Model(&model.ResourceHistory{}).
		Where("store_id = ? AND recorded_at >= ? AND recorded_at <= ?", storeID, startTime, endTime).
		Select("AVG(memory_percent)").
		Scan(&avg).Error
	return avg, err
}