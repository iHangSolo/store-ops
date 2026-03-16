package service

import (
	"store-ops-server/internal/model"
	"store-ops-server/internal/repository"
	"store-ops-server/internal/ws"

	"github.com/google/uuid"
)

type ProcessService struct {
	storeRepo *repository.StoreRepository
	auditRepo *repository.AuditLogRepository
	hub       *ws.Hub
}

func NewProcessService(
	storeRepo *repository.StoreRepository,
	auditRepo *repository.AuditLogRepository,
	hub *ws.Hub,
) *ProcessService {
	return &ProcessService{
		storeRepo: storeRepo,
		auditRepo: auditRepo,
		hub:       hub,
	}
}

type ProcessesResponse struct {
	Processes []ProcessInfo `json:"processes"`
	RequestID string        `json:"request_id"`
}

type ProcessInfo struct {
	Name          string  `json:"name"`
	PID           int32   `json:"pid"`
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryPercent float64 `json:"memory_percent"`
}

func (s *ProcessService) GetProcesses(storeID uuid.UUID) (*ProcessesResponse, error) {
	if !s.hub.IsOnline(storeID) {
		return nil, ws.ErrStoreOffline
	}

	requestID := uuid.New().String()
	msg := ws.NewMessage(ws.MessageTypeGetProcesses, map[string]interface{}{
		"request_id": requestID,
	})

	if err := s.hub.SendToStore(storeID, msg); err != nil {
		return nil, err
	}

	return &ProcessesResponse{
		Processes: []ProcessInfo{},
		RequestID: requestID,
	}, nil
}

type KillProcessResponse struct {
	RequestID string `json:"request_id"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
}

func (s *ProcessService) KillProcess(storeID uuid.UUID, pid int32, operator string) (*KillProcessResponse, error) {
	if !s.hub.IsOnline(storeID) {
		return nil, ws.ErrStoreOffline
	}

	requestID := uuid.New().String()
	msg := ws.NewMessage(ws.MessageTypeKillProcess, map[string]interface{}{
		"request_id": requestID,
		"pid":        pid,
	})

	if err := s.hub.SendToStore(storeID, msg); err != nil {
		return nil, err
	}

	s.auditRepo.Create(&model.AuditLog{
		StoreID:    &storeID,
		Operator:   operator,
		ActionType: "kill_process",
		Action:     "终止进程",
		Detail:     "终止进程 PID: " + string(pid),
	})

	return &KillProcessResponse{
		RequestID: requestID,
		Success:   true,
		Message:   "已发送终止进程请求",
	}, nil
}

type ServicesResponse struct {
	Services  []ServiceInfo `json:"services"`
	RequestID string        `json:"request_id"`
}

type ServiceInfo struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Status      string `json:"status"`
}

func (s *ProcessService) GetServices(storeID uuid.UUID) (*ServicesResponse, error) {
	if !s.hub.IsOnline(storeID) {
		return nil, ws.ErrStoreOffline
	}

	requestID := uuid.New().String()
	msg := ws.NewMessage(ws.MessageTypeGetServices, map[string]interface{}{
		"request_id": requestID,
	})

	if err := s.hub.SendToStore(storeID, msg); err != nil {
		return nil, err
	}

	return &ServicesResponse{
		Services:  []ServiceInfo{},
		RequestID: requestID,
	}, nil
}

type ControlServiceResponse struct {
	RequestID string `json:"request_id"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
}

func (s *ProcessService) ControlService(storeID uuid.UUID, serviceName, action, operator string) (*ControlServiceResponse, error) {
	if !s.hub.IsOnline(storeID) {
		return nil, ws.ErrStoreOffline
	}

	requestID := uuid.New().String()
	msg := ws.NewMessage(ws.MessageTypeControlService, map[string]interface{}{
		"request_id":   requestID,
		"service_name": serviceName,
		"action":       action,
	})

	if err := s.hub.SendToStore(storeID, msg); err != nil {
		return nil, err
	}

	s.auditRepo.Create(&model.AuditLog{
		StoreID:    &storeID,
		Operator:   operator,
		ActionType: "control_service",
		Action:     "控制服务",
		Detail:     action + " 服务: " + serviceName,
	})

	return &ControlServiceResponse{
		RequestID: requestID,
		Success:   true,
		Message:   "已发送服务控制请求",
	}, nil
}