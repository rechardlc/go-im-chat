package ctxdata

import (
	"errors"
	"context"

	"github.com/golang-jwt/jwt/v5"
)

// ParseToken 解析与验证 JWT token。
func (m *TokenManager) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("切勿尝试未知的签名方法")
		}
		return m.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token 已过期")
		}
		return nil, errors.New("无效的 token")
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token 解析失败")
}

// GetUserID 获取用户ID
func (m *TokenManager) GetUserID(tokenString string) (string, error) {
	claims, err := m.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.UserID, nil
}
// 获取JWT token
func GetJwtUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value("user_id").(string); if !ok {
		return "", errors.New("user_id 不存在")
	}
	return userID, nil
}