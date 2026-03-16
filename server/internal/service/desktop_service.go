package service

import (
	"errors"
	"time"

	"store-ops-server/internal/model"
	"store-ops-server/internal/repository"
	"store-ops-server/internal/ws"

	"github.com/google/uuid"
)

type DesktopService struct {
	desktopRepo    *repository.DesktopSessionRepository
	desktopQueue   *repository.DesktopQueueRepository
	storeRepo      *repository.StoreRepository
	auditRepo      *repository.AuditLogRepository
	hub            *ws.Hub
}

func NewDesktopService(
	desktopRepo *repository.DesktopSessionRepository,
	desktopQueue *repository.DesktopQueueRepository,
	storeRepo *repository.StoreRepository,
	auditRepo *repository.AuditLogRepository,
	hub *ws.Hub,
) *DesktopService {
	return &DesktopService{
		desktopRepo:  desktopRepo,
		desktopQueue: desktopQueue,
		storeRepo:    storeRepo,
		auditRepo:    auditRepo,
		hub:          hub,
	}
}

type DesktopStatusResponse struct {
	IsOnline      bool   `json:"is_online"`
	RustDeskID    string `json:"rustdesk_id"`
	HasActiveUser bool   `json:"has_active_user"`
	QueueLength   int    `json:"queue_length"`
}

func (s *DesktopService) GetStatus(storeID uuid.UUID) (*DesktopStatusResponse, error) {
	store, err := s.storeRepo.FindByID(storeID)
	if err != nil {
		return nil, err
	}

	queue, _ := s.desktopQueue.FindByStoreID(storeID)

	_, err = s.desktopRepo.FindActiveByStoreID(storeID)
	hasActiveUser := err == nil

	return &DesktopStatusResponse{
		IsOnline:      s.hub.IsOnline(storeID),
		RustDeskID:    store.RustDeskID,
		HasActiveUser: hasActiveUser,
		QueueLength:   len(queue),
	}, nil
}

func (s *DesktopService) Connect(storeID uuid.UUID, operator string) (*DesktopSessionResponse, error) {
	// 检查门店是否在线
	if !s.hub.IsOnline(storeID) {
		return nil, errors.New("门店离线")
	}

	// 检查是否已有活动会话
	activeSession, err := s.desktopRepo.FindActiveByStoreID(storeID)
	if err == nil {
		return &DesktopSessionResponse{
			SessionID: activeSession.ID,
			Status:    "busy",
			Message:   "已有其他用户正在连接",
		}, nil
	}

	// 检查排队
	queue, _ := s.desktopQueue.FindByStoreID(storeID)
	if len(queue) > 0 {
		return nil, errors.New("当前有用户排队，请先排队等待")
	}

	// 获取门店信息
	store, err := s.storeRepo.FindByID(storeID)
	if err != nil {
		return nil, err
	}

	// 创建会话
	session := &model.DesktopSession{
		StoreID:   storeID,
		Operator:  operator,
		StartedAt: time.Now(),
		Status:    model.DesktopSessionActive,
	}
	if err := s.desktopRepo.Create(session); err != nil {
		return nil, err
	}

	// 记录审计日志
	s.auditRepo.Create(&model.AuditLog{
		StoreID:    &storeID,
		Operator:   operator,
		ActionType: "desktop_connect",
		Action:     "远程桌面连接",
		Detail:     "开始远程桌面会话",
	})

	return &DesktopSessionResponse{
		SessionID:  session.ID,
		RustDeskID: store.RustDeskID,
		Status:     "connected",
		Message:    "连接成功",
	}, nil
}

type DesktopSessionResponse struct {
	SessionID  uuid.UUID `json:"session_id,omitempty"`
	RustDeskID string    `json:"rustdesk_id,omitempty"`
	Status     string    `json:"status"`
	Message    string    `json:"message"`
}

func (s *DesktopService) Queue(storeID uuid.UUID, operator string) (*QueueResponse, error) {
	// 检查门店是否在线
	if !s.hub.IsOnline(storeID) {
		return nil, errors.New("门店离线")
	}

	// 获取队列位置
	position, err := s.desktopQueue.GetQueuePosition(storeID)
	if err != nil {
		return nil, err
	}

	// 创建排队记录
	queue := &model.DesktopQueue{
		StoreID:  storeID,
		Operator: operator,
		Position: position,
	}
	if err := s.desktopQueue.Create(queue); err != nil {
		return nil, err
	}

	s.auditRepo.Create(&model.AuditLog{
		StoreID:    &storeID,
		Operator:   operator,
		ActionType: "desktop_queue",
		Action:     "远程桌面排队",
		Detail:     "加入远程桌面等待队列",
	})

	return &QueueResponse{
		QueueID:  queue.ID,
		Position: position,
		Message:  "已加入排队",
	}, nil
}

type QueueResponse struct {
	QueueID  uuid.UUID `json:"queue_id"`
	Position int       `json:"position"`
	Message  string    `json:"message"`
}

func (s *DesktopService) EndSession(sessionID uuid.UUID) error {
	return s.desktopRepo.EndSession(sessionID)
}

func (s *DesktopService) LeaveQueue(queueID uuid.UUID) error {
	queue, err := s.desktopQueue.FindByID(queueID)
	if err != nil {
		return err
	}

	storeID := queue.StoreID
	if err := s.desktopQueue.Delete(queueID); err != nil {
		return err
	}

	// 重新排序
	s.desktopQueue.ReorderPositions(storeID)
	return nil
}