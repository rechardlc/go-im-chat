// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"github.com/rechardlc/go-im-chat/apps/social/api/internal/config"
	"github.com/rechardlc/go-im-chat/apps/social/rpc/social_client"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/user_client"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	SocialRpc social_client.Social
	UserRpc user_client.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		SocialRpc: social_client.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
		UserRpc: user_client.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
