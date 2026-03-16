package handler

import (
	"net/http"
	"time"

	"store-ops-server/internal/pkg/response"
	"store-ops-server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuditHandler struct {
	auditService *service.AuditService
}

func NewAuditHandler(auditService *service.AuditService) *AuditHandler {
	return &AuditHandler{auditService: auditService}
}

func (h *AuditHandler) ListAuditLogs(c *gin.Context) {
	page := parseIntQuery(c, "page", 1)
	pageSize := parseIntQuery(c, "page_size", 20)

	filter := &service.AuditLogFilter{
		ActionType: c.Query("action_type"),
	}

	if storeID := c.Query("store_id"); storeID != "" {
		if id, err := uuid.Parse(storeID); err == nil {
			filter.StoreID = &id
		}
	}

	if startDate := c.Query("start_date"); startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			filter.StartDate = &t
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			filter.EndDate = &t
		}
	}

	result, err := h.auditService.ListLogs(filter, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取审计日志失败")
		return
	}

	response.Success(c, result)
}

func (h *AuditHandler) GetAuditLog(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的日志ID")
		return
	}

	result, err := h.auditService.GetLog(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "日志不存在")
		return
	}

	response.Success(c, result)
}