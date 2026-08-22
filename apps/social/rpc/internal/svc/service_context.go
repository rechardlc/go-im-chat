package svc

import (
	"github.com/rechardlc/go-im-chat/apps/social/models"
	"github.com/rechardlc/go-im-chat/apps/social/rpc/internal/config"
	"github.com/zeromicro/go-zero/core/stores/postgres"
)

type ServiceContext struct {
	Config config.Config
	models.GroupsModel
	models.FriendRequestsModel
	models.FriendsModel
	models.GroupMembersModel
	models.GroupRequestsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := postgres.New(c.Postgres.DataSource)
	return &ServiceContext{
		Config: c,
		GroupsModel: models.NewGroupsModel(sqlConn, c.CacheRedis),
		FriendRequestsModel: models.NewFriendRequestsModel(sqlConn, c.CacheRedis),
		FriendsModel: models.NewFriendsModel(sqlConn, c.CacheRedis),
		GroupMembersModel: models.NewGroupMembersModel(sqlConn, c.CacheRedis),
		GroupRequestsModel: models.NewGroupRequestsModel(sqlConn, c.CacheRedis),
	}
}
