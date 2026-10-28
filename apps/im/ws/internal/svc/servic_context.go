package svc

import (
	"github.com/rechardlc/go-im-chat/apps/im/ws/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}