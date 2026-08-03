package ctxdata

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

// ParseToken 解析与验证 JWT Token
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 校验签名算法是否匹配（防止签名算法替换攻击/None算法攻击）
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("切勿尝试未知的签名方法")
		}
		return jwtSecret, nil
	})

	if err != nil {
		// 区分具体的错误类型（如已过期）
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token 已过期")
		}
		return nil, errors.New("无效的 token")
	}

	// 校验并提取 Claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token 解析失败")
}