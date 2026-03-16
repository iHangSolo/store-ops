package repository

import (
	"time"

	"store-ops-server/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) FindByUsername(username string) (*model.Admin, error) {
	var admin model.Admin
	err := r.db.Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *AdminRepository) FindByID(id uuid.UUID) (*model.Admin, error) {
	var admin model.Admin
	err := r.db.Where("id = ?", id).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *AdminRepository) Create(admin *model.Admin) error {
	return r.db.Create(admin).Error
}

func (r *AdminRepository) UpdateLastLogin(id uuid.UUID) error {
	now := time.Now()
	return r.db.Model(&model.Admin{}).Where("id = ?", id).Update("last_login", now).Error
}