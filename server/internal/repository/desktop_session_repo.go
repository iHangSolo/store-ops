package repository

import (
	"store-ops-server/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DesktopSessionRepository struct {
	db *gorm.DB
}

func NewDesktopSessionRepository(db *gorm.DB) *DesktopSessionRepository {
	return &DesktopSessionRepository{db: db}
}

func (r *DesktopSessionRepository) Create(session *model.DesktopSession) error {
	return r.db.Create(session).Error
}

func (r *DesktopSessionRepository) Update(session *model.DesktopSession) error {
	return r.db.Save(session).Error
}

func (r *DesktopSessionRepository) FindByID(id uuid.UUID) (*model.DesktopSession, error) {
	var session model.DesktopSession
	err := r.db.Where("id = ?", id).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *DesktopSessionRepository) FindActiveByStoreID(storeID uuid.UUID) (*model.DesktopSession, error) {
	var session model.DesktopSession
	err := r.db.Where("store_id = ? AND status = ?", storeID, model.DesktopSessionActive).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *DesktopSessionRepository) EndSession(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&model.DesktopSession{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":   model.DesktopSessionCompleted,
		"ended_at": now,
	}).Error
}

func (r *DesktopSessionRepository) FindByStoreID(storeID uuid.UUID, page, pageSize int) ([]model.DesktopSession, int64, error) {
	var sessions []model.DesktopSession
	var total int64

	query := r.db.Model(&model.DesktopSession{}).Where("store_id = ?", storeID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("started_at DESC").Offset(offset).Limit(pageSize).Find(&sessions).Error
	return sessions, total, err
}

type DesktopQueueRepository struct {
	db *gorm.DB
}

func NewDesktopQueueRepository(db *gorm.DB) *DesktopQueueRepository {
	return &DesktopQueueRepository{db: db}
}

func (r *DesktopQueueRepository) Create(queue *model.DesktopQueue) error {
	return r.db.Create(queue).Error
}

func (r *DesktopQueueRepository) FindByID(id uuid.UUID) (*model.DesktopQueue, error) {
	var queue model.DesktopQueue
	err := r.db.Where("id = ?", id).First(&queue).Error
	if err != nil {
		return nil, err
	}
	return &queue, nil
}

func (r *DesktopQueueRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.DesktopQueue{}, id).Error
}

func (r *DesktopQueueRepository) FindByStoreID(storeID uuid.UUID) ([]model.DesktopQueue, error) {
	var queue []model.DesktopQueue
	err := r.db.Where("store_id = ?", storeID).Order("position ASC").Find(&queue).Error
	return queue, err
}

func (r *DesktopQueueRepository) GetQueuePosition(storeID uuid.UUID) (int, error) {
	var count int64
	err := r.db.Model(&model.DesktopQueue{}).Where("store_id = ?", storeID).Count(&count).Error
	return int(count) + 1, err
}

func (r *DesktopQueueRepository) GetFirstInQueue(storeID uuid.UUID) (*model.DesktopQueue, error) {
	var queue model.DesktopQueue
	err := r.db.Where("store_id = ?", storeID).Order("position ASC").First(&queue).Error
	if err != nil {
		return nil, err
	}
	return &queue, nil
}

func (r *DesktopQueueRepository) ReorderPositions(storeID uuid.UUID) error {
	var queue []model.DesktopQueue
	if err := r.db.Where("store_id = ?", storeID).Order("position ASC").Find(&queue).Error; err != nil {
		return err
	}

	for i, q := range queue {
		if err := r.db.Model(&model.DesktopQueue{}).Where("id = ?", q.ID).Update("position", i+1).Error; err != nil {
			return err
		}
	}
	return nil
}