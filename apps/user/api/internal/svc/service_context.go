// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"github.com/rechardlc/go-im-chat/apps/user/api/internal/config"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/user_client"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc user_client.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: user_client.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
