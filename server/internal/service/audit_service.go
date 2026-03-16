package service

import (
	"time"

	"store-ops-server/internal/model"
	"store-ops-server/internal/repository"

	"github.com/google/uuid"
)

type AuditService struct {
	auditRepo *repository.AuditLogRepository
}

func NewAuditService(auditRepo *repository.AuditLogRepository) *AuditService {
	return &AuditService{auditRepo: auditRepo}
}

type AuditLogListResponse struct {
	Logs     []model.AuditLogWithStore `json:"logs"`
	Total    int64                     `json:"total"`
	Page     int                       `json:"page"`
	PageSize int                       `json:"page_size"`
}

type AuditLogFilter struct {
	StoreID    *uuid.UUID
	ActionType string
	StartDate  *time.Time
	EndDate    *time.Time
}

func (s *AuditService) ListLogs(filter *AuditLogFilter, page, pageSize int) (*AuditLogListResponse, error) {
	logs, total, err := s.auditRepo.FindAll(
		filter.StoreID,
		filter.ActionType,
		filter.StartDate,
		filter.EndDate,
		page,
		pageSize,
	)
	if err != nil {
		return nil, err
	}

	return &AuditLogListResponse{
		Logs:     logs,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *AuditService) GetLog(id uuid.UUID) (*model.AuditLog, error) {
	return s.auditRepo.FindByID(id)
}

func (s *AuditService) CreateLog(log *model.AuditLog) error {
	return s.auditRepo.Create(log)
}

func (s *AuditService) CleanupOldLogs(days int) error {
	return s.auditRepo.DeleteOlderThan(days)
}