package ws

import (
	"time"
)

// MessageType 消息类型
type MessageType string

const (
	// 客户端 -> 服务端
	MessageTypePing           MessageType = "ping"
	MessageTypeRegister       MessageType = "register"
	MessageTypeResourceReport MessageType = "resource_report"
	MessageTypeCommandResult  MessageType = "command_result"
	MessageTypeSoftwareProgress MessageType = "software_progress"
	MessageTypeSoftwareResult   MessageType = "software_result"
	MessageTypeProcesses       MessageType = "processes"
	MessageTypeKillProcessResult MessageType = "kill_process_result"
	MessageTypeServices        MessageType = "services"
	MessageTypeControlServiceResult MessageType = "control_service_result"

	// 服务端 -> 客户端
	MessageTypePong           MessageType = "pong"
	MessageTypeRegisterAck    MessageType = "register_ack"
	MessageTypeExecuteCommand MessageType = "execute_command"
	MessageTypeExecuteScript  MessageType = "execute_script"
	MessageTypeSoftwarePush   MessageType = "software_push"
	MessageTypeGetProcesses   MessageType = "get_processes"
	MessageTypeKillProcess    MessageType = "kill_process"
	MessageTypeGetServices    MessageType = "get_services"
	MessageTypeControlService MessageType = "control_service"
	MessageTypeError          MessageType = "error"
)

// Message WebSocket 消息
type Message struct {
	Type      MessageType            `json:"type"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// RegisterPayload 注册消息负载
type RegisterPayload struct {
	DeviceID          string `json:"device_id"`
	DeviceFingerprint string `json:"device_fingerprint"`
	StoreName         string `json:"store_name"`
	RustDeskID        string `json:"rustdesk_id"`
	ClientVersion     string `json:"client_version"`
}

// RegisterAckPayload 注册确认消息负载
type RegisterAckPayload struct {
	Success  bool   `json:"success"`
	Approved bool   `json:"approved"`
	Message  string `json:"message"`
}

// ResourceReportPayload 资源上报消息负载
type ResourceReportPayload struct {
	CPUPercent        float64        `json:"cpu_percent"`
	MemoryPercent     float64        `json:"memory_percent"`
	MemoryTotalGB     float64        `json:"memory_total_gb"`
	MemoryAvailableGB float64        `json:"memory_available_gb"`
	Disks             []DiskInfoPayload `json:"disks"`
}

// DiskInfoPayload 磁盘信息
type DiskInfoPayload struct {
	Name        string  `json:"name"`
	Percent     float64 `json:"percent"`
	TotalGB     float64 `json:"total_gb"`
	AvailableGB float64 `json:"available_gb"`
}

// ExecuteCommandPayload 执行命令消息负载
type ExecuteCommandPayload struct {
	CommandID string `json:"command_id"`
	Command   string `json:"command"`
	Timeout   int    `json:"timeout"`
}

// ExecuteScriptPayload 执行脚本消息负载
type ExecuteScriptPayload struct {
	CommandID     string `json:"command_id"`
	ScriptType    string `json:"script_type"`
	ScriptContent string `json:"script_content"`
	Timeout       int    `json:"timeout"`
}

// CommandResultPayload 命令结果消息负载
type CommandResultPayload struct {
	CommandID  string `json:"command_id"`
	Status     string `json:"status"`
	Output     string `json:"output"`
	ExitCode   int    `json:"exit_code"`
	DurationMs int    `json:"duration_ms"`
}

// SoftwarePushPayload 软件推送消息负载
type SoftwarePushPayload struct {
	TaskID       string `json:"task_id"`
	DownloadURL  string `json:"download_url"`
	FileName     string `json:"file_name"`
	FileSize     int64  `json:"file_size"`
	InstallArgs  string `json:"install_args"`
}

// SoftwareProgressPayload 软件进度消息负载
type SoftwareProgressPayload struct {
	TaskID  string `json:"task_id"`
	Status  string `json:"status"`
	Progress int   `json:"progress"`
}

// SoftwareResultPayload 软件结果消息负载
type SoftwareResultPayload struct {
	TaskID  string `json:"task_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// ProcessInfoPayload 进程信息
type ProcessInfoPayload struct {
	Name          string  `json:"name"`
	PID           int32   `json:"pid"`
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryPercent float64 `json:"memory_percent"`
}

// ProcessesPayload 进程列表消息负载
type ProcessesPayload struct {
	RequestID string                `json:"request_id"`
	Processes []ProcessInfoPayload `json:"processes"`
}

// KillProcessPayload 终止进程消息负载
type KillProcessPayload struct {
	RequestID string `json:"request_id"`
	PID       int32  `json:"pid"`
}

// KillProcessResultPayload 终止进程结果消息负载
type KillProcessResultPayload struct {
	RequestID string `json:"request_id"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
}

// ServiceInfoPayload 服务信息
type ServiceInfoPayload struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Status      string `json:"status"`
}

// ServicesPayload 服务列表消息负载
type ServicesPayload struct {
	RequestID string               `json:"request_id"`
	Services  []ServiceInfoPayload `json:"services"`
}

// ControlServicePayload 控制服务消息负载
type ControlServicePayload struct {
	RequestID   string `json:"request_id"`
	ServiceName string `json:"service_name"`
	Action      string `json:"action"` // start, stop, restart
}

// ControlServiceResultPayload 控制服务结果消息负载
type ControlServiceResultPayload struct {
	RequestID string `json:"request_id"`
	Success   bool   `json:"success"`
	Message   string `json:"message"`
}

// ErrorPayload 错误消息负载
type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewMessage 创建新消息
func NewMessage(msgType MessageType, payload map[string]interface{}) *Message {
	return &Message{
		Type:      msgType,
		Payload:   payload,
		Timestamp: time.Now(),
	}
}