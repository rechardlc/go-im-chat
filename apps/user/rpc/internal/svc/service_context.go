package svc

import (
	"github.com/rechardlc/go-im-chat/apps/user/rpc/internal/config"
	"github.com/rechardlc/go-im-chat/apps/user/models"
	"github.com/zeromicro/go-zero/core/stores/postgres"
)

type ServiceContext struct {
	Config config.Config
	models.UsersModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := postgres.New(c.Postgres.DataSource)
	return &ServiceContext{
		Config:     c,
		UsersModel: models.NewUsersModel(sqlConn, c.CacheRedis),
	}
}