package middleware

import (
	"net/http"
	"strings"
	"test1/utils"

	"github.com/gin-gonic/gin"
)

// JWT 认证中间件
// TODO: 任务3 - 实现这个中间件
// 要求：
//  1. 从 Header 获取 Authorization
//  2. 格式应该是 "Bearer <token>"
//  3. 解析 Token
//  4. 把 userID 和 username 存到 Context
//  5. 失败返回 401
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 获取 Authorization Header
		authHeader := c.GetHeader("Authorization")

		// TODO: 检查是否为空
		if authHeader == "" {
			utils.ErrorWithStatus(c, http.StatusUnauthorized, utils.CodeUnauthorized, "未提供 Token")
			c.Abort()
			return
		}

		// TODO: 检查格式 "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorWithStatus(c, http.StatusUnauthorized, utils.CodeUnauthorized, "Token 格式错误")
			c.Abort()
			return
		}

		// TODO: 解析 Token
		// TODO: 把用户信息存到 Context
		// c.Set("userID", claims.UserID)
		// c.Set("username", claims.Username)

		// TODO: 调用 c.Next() 继续处理

		// 下面是示例，你需要补充完整
		claims,err:=utils.ParseToken(parts[1]) // tokenString，用这个去解析
		if err!=nil{
			utils.ErrorWithStatus(c,http.StatusUnauthorized,utils.CodeUnauthorized,"token无效")
			c.Abort()
			return 
		}
		c.Set("userID",claims.UserID)
		c.Set("username",claims.Username)
		c.Next()
	}
}
