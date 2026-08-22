package config

import "github.com/zeromicro/go-zero/zrpc"
import "time"
import "github.com/zeromicro/go-zero/core/stores/cache"

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
