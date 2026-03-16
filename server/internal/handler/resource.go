package handler

import (
	"net/http"
	"time"

	"store-ops-server/internal/pkg/response"
	"store-ops-server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ResourceHandler struct {
	resourceService *service.ResourceService
}

func NewResourceHandler(resourceService *service.ResourceService) *ResourceHandler {
	return &ResourceHandler{resourceService: resourceService}
}

func (h *ResourceHandler) GetResource(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	result, err := h.resourceService.GetCurrentResource(id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取资源信息失败")
		return
	}

	response.Success(c, result)
}

func (h *ResourceHandler) GetResourceHistory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	startTime, endTime := parseTimeRange(c)
	page := parseIntQuery(c, "page", 1)
	pageSize := parseIntQuery(c, "page_size", 100)

	result, err := h.resourceService.GetResourceHistory(id, startTime, endTime, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取资源历史失败")
		return
	}

	response.Success(c, result)
}

func (h *ResourceHandler) GetResourceStats(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	startTime, endTime := parseTimeRange(c)

	result, err := h.resourceService.GetResourceStats(id, startTime, endTime)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取资源统计失败")
		return
	}

	response.Success(c, result)
}

func parseTimeRange(c *gin.Context) (time.Time, time.Time) {
	var startTime, endTime time.Time

	if start := c.Query("start_time"); start != "" {
		if t, err := time.Parse(time.RFC3339, start); err == nil {
			startTime = t
		}
	}

	if end := c.Query("end_time"); end != "" {
		if t, err := time.Parse(time.RFC3339, end); err == nil {
			endTime = t
		}
	}

	// 默认最近24小时
	if startTime.IsZero() {
		startTime = time.Now().Add(-24 * time.Hour)
	}
	if endTime.IsZero() {
		endTime = time.Now()
	}

	return startTime, endTime
}