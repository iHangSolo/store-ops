package handler

import (
	"net/http"

	"store-ops-server/internal/model"
	"store-ops-server/internal/pkg/response"
	"store-ops-server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StoreHandler struct {
	storeService *service.StoreService
}

func NewStoreHandler(storeService *service.StoreService) *StoreHandler {
	return &StoreHandler{storeService: storeService}
}

func (h *StoreHandler) ListStores(c *gin.Context) {
	status := c.Query("status")
	page := parseIntQuery(c, "page", 1)
	pageSize := parseIntQuery(c, "page_size", 20)

	var storeStatus model.StoreStatus
	if status != "" {
		storeStatus = model.StoreStatus(status)
	}

	result, err := h.storeService.ListStores(storeStatus, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取门店列表失败")
		return
	}

	response.Success(c, result)
}

func (h *StoreHandler) GetStore(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	result, err := h.storeService.GetStore(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "门店不存在")
		return
	}

	response.Success(c, result)
}

func (h *StoreHandler) ApproveStore(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	operator := c.GetString("username")
	if err := h.storeService.ApproveStore(id, operator); err != nil {
		response.Error(c, http.StatusInternalServerError, "审批失败")
		return
	}

	response.Success(c, gin.H{"message": "审批成功"})
}

func (h *StoreHandler) RejectStore(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	operator := c.GetString("username")
	if err := h.storeService.RejectStore(id, operator, req.Reason); err != nil {
		response.Error(c, http.StatusInternalServerError, "拒绝失败")
		return
	}

	response.Success(c, gin.H{"message": "已拒绝"})
}

func (h *StoreHandler) DeleteStore(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	operator := c.GetString("username")
	if err := h.storeService.DeleteStore(id, operator); err != nil {
		response.Error(c, http.StatusInternalServerError, "删除失败")
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

func (h *StoreHandler) RegenerateToken(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	operator := c.GetString("username")
	token, err := h.storeService.RegenerateToken(id, operator)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "重新生成TOKEN失败")
		return
	}

	response.Success(c, gin.H{"token": token})
}