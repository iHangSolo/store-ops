package repository

import (
	"store-ops-server/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SoftwareTaskRepository struct {
	db *gorm.DB
}

func NewSoftwareTaskRepository(db *gorm.DB) *SoftwareTaskRepository {
	return &SoftwareTaskRepository{db: db}
}

func (r *SoftwareTaskRepository) Create(task *model.SoftwareTask) error {
	return r.db.Create(task).Error
}

func (r *SoftwareTaskRepository) Update(task *model.SoftwareTask) error {
	return r.db.Save(task).Error
}

func (r *SoftwareTaskRepository) FindByID(id uuid.UUID) (*model.SoftwareTask, error) {
	var task model.SoftwareTask
	err := r.db.Where("id = ?", id).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *SoftwareTaskRepository) FindByStoreID(storeID uuid.UUID, page, pageSize int) ([]model.SoftwareTask, int64, error) {
	var tasks []model.SoftwareTask
	var total int64

	query := r.db.Model(&model.SoftwareTask{}).Where("store_id = ?", storeID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&tasks).Error
	return tasks, total, err
}

func (r *SoftwareTaskRepository) FindPendingByStoreID(storeID uuid.UUID) ([]model.SoftwareTask, error) {
	var tasks []model.SoftwareTask
	err := r.db.Where("store_id = ? AND status = ?", storeID, model.SoftwareStatusPending).Find(&tasks).Error
	return tasks, err
}

func (r *SoftwareTaskRepository) UpdateStatus(id uuid.UUID, status model.SoftwareStatus, message string) error {
	return r.db.Model(&model.SoftwareTask{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":  status,
		"message": message,
	}).Error
}