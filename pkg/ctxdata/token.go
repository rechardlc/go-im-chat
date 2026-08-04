package ctxdata

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const tokenIssuer = "min-chat"

// CustomClaims JWT 自定义声明。
type CustomClaims struct {
	UserID string `json:"user_id"`
	Phone  string `json:"phone"`
	jwt.RegisteredClaims
}

// TokenManager 负责 JWT 签发与校验。
type TokenManager struct {
	secret             []byte
	accessTokenExpire  time.Duration
	refreshTokenExpire time.Duration
}

// NewTokenManager 根据配置创建 TokenManager。
func NewTokenManager(cfg JWTConfig) (*TokenManager, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &TokenManager{
		secret:             []byte(cfg.Secret),
		accessTokenExpire:  cfg.AccessTokenExpire,
		refreshTokenExpire: cfg.RefreshTokenExpire,
	}, nil
}

// MustNewTokenManager 创建 TokenManager，配置非法时 panic。
func MustNewTokenManager(cfg JWTConfig) *TokenManager {
	m, err := NewTokenManager(cfg)
	if err != nil {
		panic(err)
	}
	return m
}

// AccessTokenExpire 返回 access token 有效期。
func (m *TokenManager) AccessTokenExpire() time.Duration {
	return m.accessTokenExpire
}

// GenerateToken 生成 JWT access token。
func (m *TokenManager) GenerateToken(userID, phone string) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		UserID: userID,
		Phone:  phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTokenExpire)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    tokenIssuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}
