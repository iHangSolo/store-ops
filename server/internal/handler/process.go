package handler

import (
	"net/http"

	"store-ops-server/internal/pkg/response"
	"store-ops-server/internal/service"
	"store-ops-server/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProcessHandler struct {
	processService *service.ProcessService
}

func NewProcessHandler(processService *service.ProcessService) *ProcessHandler {
	return &ProcessHandler{processService: processService}
}

func (h *ProcessHandler) GetProcesses(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	result, err := h.processService.GetProcesses(id)
	if err != nil {
		if err == ws.ErrStoreOffline {
			response.Error(c, http.StatusServiceUnavailable, "门店离线")
			return
		}
		response.Error(c, http.StatusInternalServerError, "获取进程列表失败")
		return
	}

	response.Success(c, result)
}

func (h *ProcessHandler) KillProcess(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	var req struct {
		PID int32 `json:"pid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	operator := c.GetString("username")
	result, err := h.processService.KillProcess(id, req.PID, operator)
	if err != nil {
		if err == ws.ErrStoreOffline {
			response.Error(c, http.StatusServiceUnavailable, "门店离线")
			return
		}
		response.Error(c, http.StatusInternalServerError, "终止进程失败")
		return
	}

	response.Success(c, result)
}

func (h *ProcessHandler) GetServices(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	result, err := h.processService.GetServices(id)
	if err != nil {
		if err == ws.ErrStoreOffline {
			response.Error(c, http.StatusServiceUnavailable, "门店离线")
			return
		}
		response.Error(c, http.StatusInternalServerError, "获取服务列表失败")
		return
	}

	response.Success(c, result)
}

func (h *ProcessHandler) StartService(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	serviceName := c.Param("name")
	operator := c.GetString("username")

	result, err := h.processService.ControlService(id, serviceName, "start", operator)
	if err != nil {
		if err == ws.ErrStoreOffline {
			response.Error(c, http.StatusServiceUnavailable, "门店离线")
			return
		}
		response.Error(c, http.StatusInternalServerError, "启动服务失败")
		return
	}

	response.Success(c, result)
}

func (h *ProcessHandler) StopService(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	serviceName := c.Param("name")
	operator := c.GetString("username")

	result, err := h.processService.ControlService(id, serviceName, "stop", operator)
	if err != nil {
		if err == ws.ErrStoreOffline {
			response.Error(c, http.StatusServiceUnavailable, "门店离线")
			return
		}
		response.Error(c, http.StatusInternalServerError, "停止服务失败")
		return
	}

	response.Success(c, result)
}