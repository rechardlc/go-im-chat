package ctxdata

import (
	"errors"
	"fmt"
	"time"
)

// JWTConfig JWT 签发配置。Expire 只在 rpc 颁发 token 时需要；api 解析只依赖 Secret。
type JWTConfig struct {
	Secret             string
	AccessTokenExpire  time.Duration
	RefreshTokenExpire time.Duration
}

// Validate 校验签发侧 JWT 配置，启动阶段 fail fast。
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
