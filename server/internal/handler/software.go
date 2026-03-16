package handler

import (
	"net/http"

	"store-ops-server/internal/pkg/response"
	"store-ops-server/internal/service"
	"store-ops-server/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SoftwareHandler struct {
	softwareService *service.SoftwareService
}

func NewSoftwareHandler(softwareService *service.SoftwareService) *SoftwareHandler {
	return &SoftwareHandler{softwareService: softwareService}
}

func (h *SoftwareHandler) PushSoftware(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	var req service.PushSoftwareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	operator := c.GetString("username")
	result, err := h.softwareService.PushSoftware(id, operator, &req)
	if err != nil {
		if err == ws.ErrStoreOffline {
			response.Error(c, http.StatusServiceUnavailable, "门店离线")
			return
		}
		response.Error(c, http.StatusInternalServerError, "推送软件失败")
		return
	}

	response.Success(c, result)
}

func (h *SoftwareHandler) GetSoftwareStatus(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("taskId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的任务ID")
		return
	}

	result, err := h.softwareService.GetSoftwareStatus(taskID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "任务不存在")
		return
	}

	response.Success(c, result)
}

func (h *SoftwareHandler) ListSoftwareTasks(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	page := parseIntQuery(c, "page", 1)
	pageSize := parseIntQuery(c, "page_size", 20)

	result, err := h.softwareService.ListSoftwareTasks(id, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取软件任务列表失败")
		return
	}

	response.Success(c, result)
}

func (h *SoftwareHandler) CancelSoftwareTask(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("taskId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的任务ID")
		return
	}

	if err := h.softwareService.CancelSoftwareTask(taskID); err != nil {
		response.Error(c, http.StatusInternalServerError, "取消任务失败")
		return
	}

	response.Success(c, gin.H{"message": "已取消"})
}