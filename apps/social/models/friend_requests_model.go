package models

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"context"
)

var _ FriendRequestsModel = (*customFriendRequestsModel)(nil)

type (
	// FriendRequestsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFriendRequestsModel.
	FriendRequestsModel interface {
		friendRequestsModel
		FindOneByUserIdAndReqUid(ctx context.Context, userId, reqUid string) (*FriendRequests, error)
	}

	customFriendRequestsModel struct {
		*defaultFriendRequestsModel
	}
)

// NewFriendRequestsModel returns a model for the database table.
func NewFriendRequestsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FriendRequestsModel {
	return &customFriendRequestsModel{
		defaultFriendRequestsModel: newFriendRequestsModel(conn, c, opts...),
	}
}

func (m *customFriendRequestsModel) FindOneByUserIdAndReqUid(ctx context.Context, userId, reqUid string) (*FriendRequests, error) {
	query := fmt.Sprintf("select %s from %s where user_id = $1 and req_uid = $2 limit 1", friendRequestsRows, m.table)
	var resp FriendRequests
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, userId, reqUid)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}