package svc

import (
	"github.com/rechardlc/go-im-chat/apps/user/models"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/internal/config"
	"github.com/rechardlc/go-im-chat/pkg/ctxdata"
	"github.com/zeromicro/go-zero/core/stores/postgres"
)

type ServiceContext struct {
	Config config.Config
	Token  *ctxdata.TokenManager
	models.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := postgres.New(c.Postgres.DataSource)
	return &ServiceContext{
		Config: c,
		Token: ctxdata.MustNewTokenManager(ctxdata.JWTConfig{
			Secret:             c.JWT.Secret,
			AccessTokenExpire:  c.JWT.AccessTokenExpire,
			RefreshTokenExpire: c.JWT.RefreshTokenExpire,
		}),
		UsersModel: models.NewUsersModel(sqlConn, c.CacheRedis),
	}
}