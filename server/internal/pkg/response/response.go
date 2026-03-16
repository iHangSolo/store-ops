package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 错误码定义
const (
	CodeSuccess         = 0
	CodeParamError      = 1001
	CodeAuthFailed      = 1002
	CodePermissionDenied = 1003
	CodeStoreOffline    = 2001
	CodeCommandTimeout  = 2002
	CodeServerError     = 5000
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMsg 成功响应带消息
func SuccessWithMsg(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// ErrorWithData 错误响应带数据
func ErrorWithData(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// ParamError 参数错误
func ParamError(c *gin.Context, message string) {
	Error(c, CodeParamError, message)
}

// AuthFailed 认证失败
func AuthFailed(c *gin.Context, message string) {
	Error(c, CodeAuthFailed, message)
}

// PermissionDenied 权限不足
func PermissionDenied(c *gin.Context) {
	Error(c, CodePermissionDenied, "权限不足")
}

// ServerError 服务器错误
func ServerError(c *gin.Context, message string) {
	Error(c, CodeServerError, message)
}