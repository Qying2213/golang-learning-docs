package utils

import (
	"fmt"
	"test1/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 自定义 JWT Claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
// TODO: 任务2 - 实现这个函数
// 参数：userID, username
// 返回：token 字符串, error
// 要求：
//   - 过期时间设为 24 小时
//   - 使用 HS256 算法
//   - 使用 config.JWTSecret 作为密钥
func GenerateToken(userID uint, username string) (string, error) {
	// TODO: 创建 Claims
	// TODO: 创建 Token
	// TODO: 签名并返回
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodPS256, claims)
	return token.SignedString([]byte(config.JWTSecret))
}

// ParseToken 解析 JWT Token
// TODO: 任务2 - 实现这个函数
// 参数：tokenString
// 返回：*Claims, error
// 要求：
//   - 验证签名
//   - 验证过期时间
//   - 返回解析后的 Claims
func ParseToken(tokenString string) (*Claims, error) {
	// TODO: 解析 Token
	// TODO: 验证并返回 Claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	// 验证并返回 Claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
