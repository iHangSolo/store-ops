package service

import (
	"store-ops-server/internal/model"
	"store-ops-server/internal/repository"
	"store-ops-server/internal/ws"
	"store-ops-server/internal/pkg/crypto"

	"github.com/google/uuid"
)

type StoreService struct {
	storeRepo  *repository.StoreRepository
	auditRepo  *repository.AuditLogRepository
	hub        *ws.Hub
}

func NewStoreService(storeRepo *repository.StoreRepository, auditRepo *repository.AuditLogRepository, hub *ws.Hub) *StoreService {
	return &StoreService{
		storeRepo: storeRepo,
		auditRepo: auditRepo,
		hub:       hub,
	}
}

type StoreListResponse struct {
	Stores    []StoreWithStatus `json:"stores"`
	Total     int64             `json:"total"`
	Page      int               `json:"page"`
	PageSize  int               `json:"page_size"`
}

type StoreWithStatus struct {
	model.Store
	IsOnline bool `json:"is_online"`
}

func (s *StoreService) ListStores(status model.StoreStatus, page, pageSize int) (*StoreListResponse, error) {
	stores, total, err := s.storeRepo.FindAll(string(status), page, pageSize)
	if err != nil {
		return nil, err
	}

	var result []StoreWithStatus
	for _, store := range stores {
		result = append(result, StoreWithStatus{
			Store:    store,
			IsOnline: s.hub.IsOnline(store.ID),
		})
	}

	return &StoreListResponse{
		Stores:   result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *StoreService) GetStore(id uuid.UUID) (*StoreWithStatus, error) {
	store, err := s.storeRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return &StoreWithStatus{
		Store:    *store,
		IsOnline: s.hub.IsOnline(id),
	}, nil
}

func (s *StoreService) ApproveStore(id uuid.UUID, operator string) error {
	store, err := s.storeRepo.FindByID(id)
	if err != nil {
		return err
	}

	store.Status = model.StoreStatusApproved
	if err := s.storeRepo.Update(store); err != nil {
		return err
	}

	s.auditRepo.Create(&model.AuditLog{
		StoreID:    &id,
		Operator:   operator,
		ActionType: "store_approve",
		Action:     "审批门店",
		Detail:     "门店审批通过: " + store.Name,
	})

	return nil
}

func (s *StoreService) RejectStore(id uuid.UUID, operator string, reason string) error {
	store, err := s.storeRepo.FindByID(id)
	if err != nil {
		return err
	}

	store.Status = model.StoreStatusRejected
	if err := s.storeRepo.Update(store); err != nil {
		return err
	}

	s.auditRepo.Create(&model.AuditLog{
		StoreID:    &id,
		Operator:   operator,
		ActionType: "store_reject",
		Action:     "拒绝门店",
		Detail:     "门店被拒绝: " + store.Name + ", 原因: " + reason,
	})

	return nil
}

func (s *StoreService) DeleteStore(id uuid.UUID, operator string) error {
	store, err := s.storeRepo.FindByID(id)
	if err != nil {
		return err
	}

	if err := s.storeRepo.Delete(id); err != nil {
		return err
	}

	s.auditRepo.Create(&model.AuditLog{
		StoreID:    &id,
		Operator:   operator,
		ActionType: "store_delete",
		Action:     "删除门店",
		Detail:     "删除门店: " + store.Name,
	})

	return nil
}

func (s *StoreService) RegenerateToken(id uuid.UUID, operator string) (string, error) {
	store, err := s.storeRepo.FindByID(id)
	if err != nil {
		return "", err
	}

	newToken := crypto.GenerateToken()
	store.Token = newToken
	if err := s.storeRepo.Update(store); err != nil {
		return "", err
	}

	s.auditRepo.Create(&model.AuditLog{
		StoreID:    &id,
		Operator:   operator,
		ActionType: "store_regen_token",
		Action:     "重新生成TOKEN",
		Detail:     "重新生成门店TOKEN: " + store.Name,
	})

	return newToken, nil
}

func (s *StoreService) GetOnlineCount() int {
	return s.hub.GetOnlineCount()
}