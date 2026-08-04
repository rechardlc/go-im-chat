package config

import (
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Postgres struct {
		DataSource string
	}
	CacheRedis cache.CacheConf
	JWT struct {
		Secret string
		AccessTokenExpire time.Duration
		RefreshTokenExpire time.Duration
	}
}
