package repository

import (
	"store-ops-server/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StoreRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) *StoreRepository {
	return &StoreRepository{db: db}
}

func (r *StoreRepository) FindAll(status string, page, pageSize int) ([]model.Store, int64, error) {
	var stores []model.Store
	var total int64

	query := r.db.Model(&model.Store{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&stores).Error
	return stores, total, err
}

func (r *StoreRepository) FindByID(id uuid.UUID) (*model.Store, error) {
	var store model.Store
	err := r.db.Where("id = ?", id).First(&store).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *StoreRepository) FindByDeviceID(deviceID uuid.UUID) (*model.Store, error) {
	var store model.Store
	err := r.db.Where("device_id = ?", deviceID).First(&store).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *StoreRepository) FindPending() ([]model.Store, error) {
	var stores []model.Store
	err := r.db.Where("approved = ?", false).Find(&stores).Error
	return stores, err
}

func (r *StoreRepository) Create(store *model.Store) error {
	return r.db.Create(store).Error
}

func (r *StoreRepository) Update(store *model.Store) error {
	return r.db.Save(store).Error
}

func (r *StoreRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Store{}, id).Error
}

func (r *StoreRepository) UpdateStatus(id uuid.UUID, status model.StoreStatus) error {
	return r.db.Model(&model.Store{}).Where("id = ?", id).Update("status", status).Error
}

func (r *StoreRepository) Approve(id uuid.UUID) error {
	return r.db.Model(&model.Store{}).Where("id = ?", id).Update("approved", true).Error
}

func (r *StoreRepository) CountOnline() (int64, error) {
	var count int64
	err := r.db.Model(&model.Store{}).Where("is_online = ?", true).Count(&count).Error
	return count, err
}

func (r *StoreRepository) FindByToken(token string) (*model.Store, error) {
	var store model.Store
	err := r.db.Where("token = ?", token).First(&store).Error
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *StoreRepository) UpdateOnlineStatus(id uuid.UUID, isOnline bool) error {
	return r.db.Model(&model.Store{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_online": isOnline,
	}).Error
}

func (r *StoreRepository) FindAllWithFilter(status model.StoreStatus, page, pageSize int) ([]model.Store, int64, error) {
	var stores []model.Store
	var total int64

	query := r.db.Model(&model.Store{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&stores).Error
	return stores, total, err
}