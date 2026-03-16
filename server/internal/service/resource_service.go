package service

import (
	"time"

	"store-ops-server/internal/model"
	"store-ops-server/internal/repository"
	"store-ops-server/internal/ws"

	"github.com/google/uuid"
)

type ResourceService struct {
	resourceRepo *repository.ResourceHistoryRepository
	storeRepo    *repository.StoreRepository
	hub          *ws.Hub
}

func NewResourceService(
	resourceRepo *repository.ResourceHistoryRepository,
	storeRepo *repository.StoreRepository,
	hub *ws.Hub,
) *ResourceService {
	return &ResourceService{
		resourceRepo: resourceRepo,
		storeRepo:    storeRepo,
		hub:          hub,
	}
}

type ResourceResponse struct {
	CPUPercent        float64        `json:"cpu_percent"`
	MemoryPercent     float64        `json:"memory_percent"`
	MemoryTotalGB     float64        `json:"memory_total_gb"`
	MemoryAvailableGB float64        `json:"memory_available_gb"`
	Disks             []model.DiskInfo `json:"disks"`
	RecordedAt        time.Time      `json:"recorded_at"`
	IsOnline          bool           `json:"is_online"`
}

func (s *ResourceService) GetCurrentResource(storeID uuid.UUID) (*ResourceResponse, error) {
	isOnline := s.hub.IsOnline(storeID)

	store, err := s.storeRepo.FindByID(storeID)
	if err != nil {
		return nil, err
	}

	// 获取最新资源历史
	latest, err := s.resourceRepo.FindLatestByStoreID(storeID)
	if err != nil {
		// 没有历史数据，返回门店缓存的资源信息
		return &ResourceResponse{
			CPUPercent:    store.CPUPercent,
			MemoryPercent: store.MemoryPercent,
			IsOnline:      isOnline,
			RecordedAt:    time.Now(),
		}, nil
	}

	return &ResourceResponse{
		CPUPercent:        latest.CPUPercent,
		MemoryPercent:     latest.MemoryPercent,
		MemoryAvailableGB: latest.MemoryAvailableGB,
		Disks:             latest.DiskData,
		RecordedAt:        latest.RecordedAt,
		IsOnline:          isOnline,
	}, nil
}

type ResourceHistoryResponse struct {
	History  []model.ResourceHistory `json:"history"`
	Total    int64                   `json:"total"`
	Page     int                     `json:"page"`
	PageSize int                     `json:"page_size"`
}

func (s *ResourceService) GetResourceHistory(storeID uuid.UUID, startTime, endTime time.Time, page, pageSize int) (*ResourceHistoryResponse, error) {
	history, total, err := s.resourceRepo.FindByStoreID(storeID, startTime, endTime, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &ResourceHistoryResponse{
		History:  history,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

type ResourceStatsResponse struct {
	AvgCPU    float64 `json:"avg_cpu"`
	AvgMemory float64 `json:"avg_memory"`
}

func (s *ResourceService) GetResourceStats(storeID uuid.UUID, startTime, endTime time.Time) (*ResourceStatsResponse, error) {
	avgCPU, err := s.resourceRepo.GetAverageCPU(storeID, startTime, endTime)
	if err != nil {
		avgCPU = 0
	}

	avgMemory, err := s.resourceRepo.GetAverageMemory(storeID, startTime, endTime)
	if err != nil {
		avgMemory = 0
	}

	return &ResourceStatsResponse{
		AvgCPU:    avgCPU,
		AvgMemory: avgMemory,
	}, nil
}