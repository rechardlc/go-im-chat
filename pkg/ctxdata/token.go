package ctxdata

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 1. 定义自定义声明结构体（Custom Claims）
// jwt.RegisteredClaims 包含了标准字段：Issuer, Subject, Audience, ExpiresAt, NotBefore, IssuedAt, ID
type CustomClaims struct {
	UserID   int64  `json:"user_id"`
	Phone    string `json:"phone"`
	jwt.RegisteredClaims
}

// 2. 密钥（Secret Key），生产环境中建议从环境变量或配置文件（如 config.yaml）读取
var jwtSecret = []byte("your-256-bit-secret-key-change-it-in-prod")

// CustomClaims 过期时间设置（例如 2 小时）
const TokenExpireDuration = 2 * time.Hour

// GenerateToken 生成 JWT Token
func GenerateToken(userID int64, phone string) (string, error) {
	// 创建自己的 Claims
	claims := CustomClaims{
		UserID:   userID,
		Phone:    phone,
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			// 签发时间
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// 生效时间
			NotBefore: jwt.NewNumericDate(time.Now()),
			// 签发人
			Issuer: "min-chat",
		},
	}

	// 使用指定的签名方法（如 HMAC-SHA256）和 Claim 创建 token 对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥进行签名并获得完整的编码后的字符串 Token
	return token.SignedString(jwtSecret)
}