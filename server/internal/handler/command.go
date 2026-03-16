package handler

import (
	"net/http"

	"store-ops-server/internal/pkg/response"
	"store-ops-server/internal/service"
	"store-ops-server/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommandHandler struct {
	commandService *service.CommandService
}

func NewCommandHandler(commandService *service.CommandService) *CommandHandler {
	return &CommandHandler{commandService: commandService}
}

func (h *CommandHandler) ExecuteCommand(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的门店ID")
		return
	}

	var req service.ExecuteCommandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	operator := c.GetString("username")
	result, err := h.commandService.ExecuteCommand(id, operator, &req)
	if err != nil {
		if err == ws.ErrStoreOffline {
			response.Error(c, http.StatusServiceUnavailable, "门店离线")
			return
		}
		response.Error(c, http.StatusInternalServerError, "执行命令失败")
		return
	}

	response.Success(c, result)
}

func (h *CommandHandler) GetCommandResult(c *gin.Context) {
	commandID, err := uuid.Parse(c.Param("commandId"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的命令ID")
		return
	}

	result, err := h.commandService.GetCommandResult(commandID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "命令不存在")
		return
	}

	response.Success(c, result)
}

func (h *CommandHandler) ListCommands(c *gin.Context) {
	page := parseIntQuery(c, "page", 1)
	pageSize := parseIntQuery(c, "page_size", 20)

	var storeID *uuid.UUID
	if id := c.Query("store_id"); id != "" {
		parsed, err := uuid.Parse(id)
		if err == nil {
			storeID = &parsed
		}
	}

	operator := c.Query("operator")

	result, err := h.commandService.ListCommands(storeID, operator, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取命令列表失败")
		return
	}

	response.Success(c, result)
}