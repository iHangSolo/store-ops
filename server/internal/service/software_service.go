package service

import (
	"store-ops-server/internal/model"
	"store-ops-server/internal/repository"
	"store-ops-server/internal/ws"

	"github.com/google/uuid"
)

type SoftwareService struct {
	softwareRepo *repository.SoftwareTaskRepository
	storeRepo    *repository.StoreRepository
	auditRepo    *repository.AuditLogRepository
	hub          *ws.Hub
}

func NewSoftwareService(
	softwareRepo *repository.SoftwareTaskRepository,
	storeRepo *repository.StoreRepository,
	auditRepo *repository.AuditLogRepository,
	hub *ws.Hub,
) *SoftwareService {
	return &SoftwareService{
		softwareRepo: softwareRepo,
		storeRepo:    storeRepo,
		auditRepo:    auditRepo,
		hub:          hub,
	}
}

type PushSoftwareRequest struct {
	FileName    string `json:"file_name" binding:"required"`
	FileSize    int64  `json:"file_size"`
	DownloadURL string `json:"download_url" binding:"required"`
	InstallArgs string `json:"install_args"`
}

type PushSoftwareResponse struct {
	TaskID string `json:"task_id"`
	Status string `json:"status"`
}

func (s *SoftwareService) PushSoftware(storeID uuid.UUID, operator string, req *PushSoftwareRequest) (*PushSoftwareResponse, error) {
	// 检查门店是否在线
	if !s.hub.IsOnline(storeID) {
		return nil, ws.ErrStoreOffline
	}

	// 创建软件任务
	taskID := uuid.New()
	task := &model.SoftwareTask{
		ID:          taskID,
		StoreID:     storeID,
		Operator:    operator,
		FileName:    req.FileName,
		FileSize:    req.FileSize,
		InstallArgs: req.InstallArgs,
		Status:      model.SoftwareStatusPending,
	}
	if err := s.softwareRepo.Create(task); err != nil {
		return nil, err
	}

	// 发送推送消息到客户端
	msg := ws.NewMessage(ws.MessageTypeSoftwarePush, map[string]interface{}{
		"task_id":      taskID.String(),
		"download_url": req.DownloadURL,
		"file_name":    req.FileName,
		"file_size":    req.FileSize,
		"install_args": req.InstallArgs,
	})

	if err := s.hub.SendToStore(storeID, msg); err != nil {
		return nil, err
	}

	// 更新任务状态
	s.softwareRepo.UpdateStatus(taskID, model.SoftwareStatusDownloading, "正在下载")

	// 记录审计日志
	s.auditRepo.Create(&model.AuditLog{
		StoreID:    &storeID,
		Operator:   operator,
		ActionType: "push_software",
		Action:     "推送软件",
		Detail:     "推送软件: " + req.FileName,
	})

	return &PushSoftwareResponse{
		TaskID: taskID.String(),
		Status: string(model.SoftwareStatusDownloading),
	}, nil
}

func (s *SoftwareService) GetSoftwareStatus(taskID uuid.UUID) (*model.SoftwareTask, error) {
	return s.softwareRepo.FindByID(taskID)
}

type SoftwareTaskListResponse struct {
	Tasks    []model.SoftwareTask `json:"tasks"`
	Total    int64                `json:"total"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"page_size"`
}

func (s *SoftwareService) ListSoftwareTasks(storeID uuid.UUID, page, pageSize int) (*SoftwareTaskListResponse, error) {
	tasks, total, err := s.softwareRepo.FindByStoreID(storeID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &SoftwareTaskListResponse{
		Tasks:    tasks,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

func (s *SoftwareService) CancelSoftwareTask(taskID uuid.UUID) error {
	task, err := s.softwareRepo.FindByID(taskID)
	if err != nil {
		return err
	}

	if task.Status == model.SoftwareStatusInstalled || task.Status == model.SoftwareStatusFailed {
		return nil // 已完成，无需取消
	}

	return s.softwareRepo.UpdateStatus(taskID, model.SoftwareStatusFailed, "用户取消")
}