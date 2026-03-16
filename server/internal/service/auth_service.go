package service

import (
	"errors"
	"time"

	"store-ops-server/internal/model"
	"store-ops-server/internal/pkg/crypto"
	"store-ops-server/internal/pkg/auth"
	"store-ops-server/internal/repository"

	"github.com/google/uuid"
)

type AuthService struct {
	adminRepo  *repository.AdminRepository
	auditRepo  *repository.AuditLogRepository
}

func NewAuthService(adminRepo *repository.AdminRepository, auditRepo *repository.AuditLogRepository) *AuthService {
	return &AuthService{
		adminRepo: adminRepo,
		auditRepo: auditRepo,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Admin     AdminInfo `json:"admin"`
}

type AdminInfo struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Nickname string    `json:"nickname"`
	Role     string    `json:"role"`
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	admin, err := s.adminRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if !crypto.CheckPassword(req.Password, admin.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	if admin.Status != model.AdminStatusActive {
		return nil, errors.New("账户已被禁用")
	}

	// 生成 JWT token
	token, expiresAt, err := auth.GenerateToken(admin.ID, admin.Username, admin.Role)
	if err != nil {
		return nil, err
	}

	// 更新最后登录时间
	s.adminRepo.UpdateLastLogin(admin.ID)

	// 记录审计日志
	s.auditRepo.Create(&model.AuditLog{
		Operator:   admin.Username,
		ActionType: "login",
		Action:     "管理员登录",
		Detail:     "登录成功",
	})

	return &LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		Admin: AdminInfo{
			ID:       admin.ID,
			Username: admin.Username,
			Nickname: admin.Nickname,
			Role:     admin.Role,
		},
	}, nil
}

func (s *AuthService) Logout(adminID uuid.UUID, username string) error {
	s.auditRepo.Create(&model.AuditLog{
		StoreID:    nil,
		Operator:   username,
		ActionType: "logout",
		Action:     "管理员登出",
		Detail:     "登出成功",
	})
	return nil
}