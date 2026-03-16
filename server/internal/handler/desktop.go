package handler

import (
	"net/http"

	"store-ops-server/internal/pkg/response"
	"store-ops-server/internal/service"
	"store-ops-server/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DesktopHandler struct {
	desktopService *service.DesktopService
}

func NewDesktopHandler(desktopService *service.DesktopService) *DesktopHandler {
	return &DesktopHandler{desktopService: desktopService}
}

func (h *DesktopHandler) GetDesktopStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	result, err := h.desktopService.GetStatus(id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取状态失败")
		return
	}

	response.Success(c, result)
}

func (h *DesktopHandler) ConnectDesktop(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	operator := c.GetString("username")
	result, err := h.desktopService.Connect(id, operator)
	if err != nil {
		if err == ws.ErrStoreOffline {
			response.Error(c, http.StatusServiceUnavailable, "门店离线")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *DesktopHandler) QueueDesktop(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	operator := c.GetString("username")
	result, err := h.desktopService.Queue(id, operator)
	if err != nil {
		if err == ws.ErrStoreOffline {
			response.Error(c, http.StatusServiceUnavailable, "门店离线")
			return
		}
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *DesktopHandler) EndSession(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("session_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的会话ID")
		return
	}

	if err := h.desktopService.EndSession(sessionID); err != nil {
		response.Error(c, http.StatusInternalServerError, "结束会话失败")
		return
	}

	response.Success(c, gin.H{"message": "会话已结束"})
}

func (h *DesktopHandler) LeaveQueue(c *gin.Context) {
	queueID, err := uuid.Parse(c.Param("queue_id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的排队ID")
		return
	}

	if err := h.desktopService.LeaveQueue(queueID); err != nil {
		response.Error(c, http.StatusInternalServerError, "离开排队失败")
		return
	}

	response.Success(c, gin.H{"message": "已离开排队"})
}