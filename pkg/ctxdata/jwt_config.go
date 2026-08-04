package ctxdata

import (
	"errors"
	"fmt"
	"time"
)

// JWTConfig JWT 签发与校验配置（由 user.yaml → config.Config 注入）。
type JWTConfig struct {
	Secret             string
	AccessTokenExpire  time.Duration
	RefreshTokenExpire time.Duration
}

// Validate 校验 JWT 配置，启动阶段 fail fast。
func (c JWTConfig) Validate() error {
	if c.Secret == "" {
		return errors.New("jwt secret is required")
	}
	if c.AccessTokenExpire <= 0 {
		return fmt.Errorf("jwt access token expire must be positive, got %v", c.AccessTokenExpire)
	}
	if c.RefreshTokenExpire <= 0 {
		return fmt.Errorf("jwt refresh token expire must be positive, got %v", c.RefreshTokenExpire)
	}
	return nil
}
