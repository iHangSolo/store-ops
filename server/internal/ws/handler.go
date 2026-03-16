package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"store-ops-server/internal/model"
	"store-ops-server/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Handler WebSocket 处理器
type Handler struct {
	hub           *Hub
	storeRepo     *repository.StoreRepository
	auditLogRepo  *repository.AuditLogRepository
	resourceRepo  *repository.ResourceHistoryRepository
	commandRepo   *repository.CommandLogRepository
	softwareRepo  *repository.SoftwareTaskRepository
	desktopRepo   *repository.DesktopSessionRepository
	desktopQueue  *repository.DesktopQueueRepository
}

// NewHandler 创建 WebSocket 处理器
func NewHandler(
	hub *Hub,
	storeRepo *repository.StoreRepository,
	auditLogRepo *repository.AuditLogRepository,
	resourceRepo *repository.ResourceHistoryRepository,
	commandRepo *repository.CommandLogRepository,
	softwareRepo *repository.SoftwareTaskRepository,
	desktopRepo *repository.DesktopSessionRepository,
	desktopQueue *repository.DesktopQueueRepository,
) *Handler {
	return &Handler{
		hub:          hub,
		storeRepo:    storeRepo,
		auditLogRepo: auditLogRepo,
		resourceRepo: resourceRepo,
		commandRepo:  commandRepo,
		softwareRepo: softwareRepo,
		desktopRepo:  desktopRepo,
		desktopQueue: desktopQueue,
	}
}

// HandleWebSocket 处理 WebSocket 连接
func (h *Handler) HandleWebSocket(c *gin.Context) {
	// 获取 TOKEN 和设备指纹
	token := c.Query("token")
	deviceFingerprint := c.Query("fingerprint")

	if token == "" || deviceFingerprint == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少认证参数"})
		return
	}

	// 验证 TOKEN 并获取门店
	store, err := h.storeRepo.FindByToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的 TOKEN"})
		return
	}

	// 验证设备指纹
	if store.DeviceFingerprint != "" && store.DeviceFingerprint != deviceFingerprint {
		c.JSON(http.StatusForbidden, gin.H{"error": "设备指纹不匹配"})
		return
	}

	// 升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket 升级失败: %v", err)
		return
	}

	// 创建客户端
	client := &Client{
		ID:               uuid.New(),
		Conn:             conn,
		StoreID:          store.ID,
		DeviceID:         store.DeviceID,
		DeviceFingerprint: deviceFingerprint,
		Send:             make(chan []byte, 256),
	}

	// 注册客户端
	h.hub.Register <- client

	// 更新门店在线状态
	h.storeRepo.UpdateOnlineStatus(store.ID, true)

	// 记录审计日志
	h.auditLogRepo.Create(&model.AuditLog{
		StoreID:    &store.ID,
		ActionType: "client_connect",
		Action:     "客户端连接",
		Detail:     "门店客户端成功连接",
	})

	// 启动读写协程
	go client.WritePump()
	go client.ReadPump(func(msg *Message) {
		h.handleMessage(client, msg)
	})
}

// handleMessage 处理客户端消息
func (h *Handler) handleMessage(client *Client, msg *Message) {
	switch msg.Type {
	case MessageTypePing:
		h.handlePing(client)

	case MessageTypeRegister:
		h.handleRegister(client, msg)

	case MessageTypeResourceReport:
		h.handleResourceReport(client, msg)

	case MessageTypeCommandResult:
		h.handleCommandResult(client, msg)

	case MessageTypeSoftwareProgress:
		h.handleSoftwareProgress(client, msg)

	case MessageTypeSoftwareResult:
		h.handleSoftwareResult(client, msg)

	case MessageTypeProcesses:
		h.handleProcesses(client, msg)

	case MessageTypeKillProcessResult:
		h.handleKillProcessResult(client, msg)

	case MessageTypeServices:
		h.handleServices(client, msg)

	case MessageTypeControlServiceResult:
		h.handleControlServiceResult(client, msg)
	}
}

// handlePing 处理心跳
func (h *Handler) handlePing(client *Client) {
	pong := NewMessage(MessageTypePong, nil)
	data, _ := json.Marshal(pong)
	client.Send <- data
}

// handleRegister 处理客户端注册
func (h *Handler) handleRegister(client *Client, msg *Message) {
	var payload RegisterPayload
	if p, ok := msg.Payload["device_id"].(string); ok {
		payload.DeviceID = p
	}
	if p, ok := msg.Payload["device_fingerprint"].(string); ok {
		payload.DeviceFingerprint = p
	}
	if p, ok := msg.Payload["store_name"].(string); ok {
		payload.StoreName = p
	}
	if p, ok := msg.Payload["rustdesk_id"].(string); ok {
		payload.RustDeskID = p
	}
	if p, ok := msg.Payload["client_version"].(string); ok {
		payload.ClientVersion = p
	}

	// 更新门店信息
	store, _ := h.storeRepo.FindByID(client.StoreID)
	if store != nil {
		if payload.RustDeskID != "" {
			store.RustDeskID = payload.RustDeskID
		}
		if payload.ClientVersion != "" {
			store.ClientVersion = payload.ClientVersion
		}
		if store.DeviceFingerprint == "" && payload.DeviceFingerprint != "" {
			store.DeviceFingerprint = payload.DeviceFingerprint
		}
		h.storeRepo.Update(store)
	}

	// 发送注册确认
	ack := NewMessage(MessageTypeRegisterAck, map[string]interface{}{
		"success":  true,
		"approved": store != nil && store.Status == model.StoreStatusApproved,
		"message":  "注册成功",
	})
	data, _ := json.Marshal(ack)
	client.Send <- data
}

// handleResourceReport 处理资源上报
func (h *Handler) handleResourceReport(client *Client, msg *Message) {
	var payload ResourceReportPayload
	payloadBytes, _ := json.Marshal(msg.Payload)
	json.Unmarshal(payloadBytes, &payload)

	// 转换磁盘信息
	var disks []model.DiskInfo
	for _, d := range payload.Disks {
		disks = append(disks, model.DiskInfo{
			Name:        d.Name,
			Percent:     d.Percent,
			TotalGB:     d.TotalGB,
			AvailableGB: d.AvailableGB,
		})
	}

	// 保存资源历史
	history := &model.ResourceHistory{
		StoreID:           client.StoreID,
		CPUPercent:        payload.CPUPercent,
		MemoryPercent:     payload.MemoryPercent,
		MemoryAvailableGB: payload.MemoryAvailableGB,
		DiskData:          disks,
		RecordedAt:        time.Now(),
	}
	h.resourceRepo.Create(history)

	// 更新门店资源信息
	store, _ := h.storeRepo.FindByID(client.StoreID)
	if store != nil {
		store.CPUPercent = payload.CPUPercent
		store.MemoryPercent = payload.MemoryPercent
		h.storeRepo.Update(store)
	}
}

// handleCommandResult 处理命令结果
func (h *Handler) handleCommandResult(client *Client, msg *Message) {
	var payload CommandResultPayload
	payloadBytes, _ := json.Marshal(msg.Payload)
	json.Unmarshal(payloadBytes, &payload)

	// 更新命令日志
	log, err := h.commandRepo.FindByID(uuid.MustParse(payload.CommandID))
	if err == nil {
		log.Output = payload.Output
		log.Status = model.CommandStatus(payload.Status)
		log.DurationMs = payload.DurationMs
		h.commandRepo.Update(log)
	}
}

// handleSoftwareProgress 处理软件进度
func (h *Handler) handleSoftwareProgress(client *Client, msg *Message) {
	var payload SoftwareProgressPayload
	payloadBytes, _ := json.Marshal(msg.Payload)
	json.Unmarshal(payloadBytes, &payload)

	// 更新软件任务状态
	task, err := h.softwareRepo.FindByID(uuid.MustParse(payload.TaskID))
	if err == nil {
		task.Status = model.SoftwareStatus(payload.Status)
		h.softwareRepo.Update(task)
	}
}

// handleSoftwareResult 处理软件结果
func (h *Handler) handleSoftwareResult(client *Client, msg *Message) {
	var payload SoftwareResultPayload
	payloadBytes, _ := json.Marshal(msg.Payload)
	json.Unmarshal(payloadBytes, &payload)

	// 更新软件任务状态
	task, err := h.softwareRepo.FindByID(uuid.MustParse(payload.TaskID))
	if err == nil {
		task.Status = model.SoftwareStatus(payload.Status)
		task.Message = payload.Message
		now := time.Now()
		task.CompletedAt = &now
		h.softwareRepo.Update(task)
	}
}

// handleProcesses 处理进程列表（存储在客户端上下文中）
func (h *Handler) handleProcesses(client *Client, msg *Message) {
	// 进程列表用于实时查询，不需要持久化
	// 可以通过 request_id 关联回原始请求
}

// handleKillProcessResult 处理终止进程结果
func (h *Handler) handleKillProcessResult(client *Client, msg *Message) {
	// 记录审计日志
	var payload KillProcessResultPayload
	payloadBytes, _ := json.Marshal(msg.Payload)
	json.Unmarshal(payloadBytes, &payload)

	h.auditLogRepo.Create(&model.AuditLog{
		StoreID:    &client.StoreID,
		ActionType: "kill_process",
		Action:     "终止进程",
		Detail:     payload.Message,
	})
}

// handleServices 处理服务列表
func (h *Handler) handleServices(client *Client, msg *Message) {
	// 服务列表用于实时查询，不需要持久化
}

// handleControlServiceResult 处理服务控制结果
func (h *Handler) handleControlServiceResult(client *Client, msg *Message) {
	var payload ControlServiceResultPayload
	payloadBytes, _ := json.Marshal(msg.Payload)
	json.Unmarshal(payloadBytes, &payload)

	h.auditLogRepo.Create(&model.AuditLog{
		StoreID:    &client.StoreID,
		ActionType: "control_service",
		Action:     "控制服务",
		Detail:     payload.Message,
	})
}