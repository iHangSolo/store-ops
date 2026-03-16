package handler

import (
	"net/http"

	"store-ops-server/internal/pkg/response"
	"store-ops-server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	result, err := h.authService.Login(&req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	adminIDStr := c.GetString("admin_id")
	username := c.GetString("username")

	adminID, err := uuid.Parse(adminIDStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的管理员ID")
		return
	}

	if err := h.authService.Logout(adminID, username); err != nil {
		response.Error(c, http.StatusInternalServerError, "登出失败")
		return
	}

	response.Success(c, gin.H{"message": "登出成功"})
}