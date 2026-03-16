package service

import (
	"time"

	"store-ops-server/internal/model"
	"store-ops-server/internal/repository"
	"store-ops-server/internal/ws"

	"github.com/google/uuid"
)

type CommandService struct {
	commandRepo *repository.CommandLogRepository
	storeRepo   *repository.StoreRepository
	auditRepo   *repository.AuditLogRepository
	hub         *ws.Hub
}

func NewCommandService(
	commandRepo *repository.CommandLogRepository,
	storeRepo *repository.StoreRepository,
	auditRepo *repository.AuditLogRepository,
	hub *ws.Hub,
) *CommandService {
	return &CommandService{
		commandRepo: commandRepo,
		storeRepo:   storeRepo,
		auditRepo:   auditRepo,
		hub:         hub,
	}
}

type ExecuteCommandRequest struct {
	Command string `json:"command" binding:"required"`
	Timeout int    `json:"timeout"`
}

type ExecuteCommandResponse struct {
	CommandID string `json:"command_id"`
	Status    string `json:"status"`
}

func (s *CommandService) ExecuteCommand(storeID uuid.UUID, operator string, req *ExecuteCommandRequest) (*ExecuteCommandResponse, error) {
	// 检查门店是否在线
	if !s.hub.IsOnline(storeID) {
		return nil, ws.ErrStoreOffline
	}

	// 创建命令日志
	commandID := uuid.New()
	commandLog := &model.CommandLog{
		ID:         commandID,
		StoreID:    &storeID,
		Operator:   operator,
		Command:    req.Command,
		Status:     model.CommandStatusExecuting,
		ExecutedAt: time.Now(),
	}
	if err := s.commandRepo.Create(commandLog); err != nil {
		return nil, err
	}

	// 发送命令到客户端
	msg := ws.NewMessage(ws.MessageTypeExecuteCommand, map[string]interface{}{
		"command_id": commandID.String(),
		"command":    req.Command,
		"timeout":    req.Timeout,
	})

	if err := s.hub.SendToStore(storeID, msg); err != nil {
		return nil, err
	}

	// 记录审计日志
	s.auditRepo.Create(&model.AuditLog{
		StoreID:    &storeID,
		Operator:   operator,
		ActionType: "execute_command",
		Action:     "执行命令",
		Detail:     req.Command,
	})

	return &ExecuteCommandResponse{
		CommandID: commandID.String(),
		Status:    string(model.CommandStatusExecuting),
	}, nil
}

func (s *CommandService) GetCommandResult(commandID uuid.UUID) (*model.CommandLog, error) {
	return s.commandRepo.FindByID(commandID)
}

type CommandListResponse struct {
	Commands []model.CommandLog `json:"commands"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

func (s *CommandService) ListCommands(storeID *uuid.UUID, operator string, page, pageSize int) (*CommandListResponse, error) {
	var commands []model.CommandLog
	var total int64
	var err error

	if storeID != nil {
		commands, total, err = s.commandRepo.FindByStoreID(*storeID, page, pageSize)
	} else if operator != "" {
		commands, total, err = s.commandRepo.FindByOperator(operator, page, pageSize)
	} else {
		// 返回空列表
		commands = []model.CommandLog{}
		total = 0
	}

	if err != nil {
		return nil, err
	}

	return &CommandListResponse{
		Commands: commands,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}