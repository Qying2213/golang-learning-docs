package utils

import (
	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应带消息
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(200, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(200, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// ErrorWithStatus 错误响应（带 HTTP 状态码）
func ErrorWithStatus(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// 常用错误码
const (
	CodeSuccess       = 0
	CodeParamError    = 40001
	CodeUnauthorized  = 40101
	CodeForbidden     = 40301
	CodeNotFound      = 40401
	CodeUserExist     = 40002
	CodeUserNotExist  = 40003
	CodePasswordError = 40004
	CodeServerError   = 50001
)
