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
		Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error
		UpdateWithSession(ctx context.Context, session sqlx.Session, data *FriendRequests) error
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
/**
* 根据用户ID和请求用户ID查找好友申请记录
* @param ctx context.Context
* @param userId string
* @param reqUid string
* @return *FriendRequests, error
*/
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

/**
* 事务处理
* @param ctx context.Context
* @param fn func(ctx context.Context, session sqlx.Session) error
* @return error
*/
func (m *customFriendRequestsModel) Trans(ctx context.Context, fn func(ctx context.Context, session sqlx.Session) error) error {
	return m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		return fn(ctx, session)
	})
}

/**
* 更新好友申请记录
* @param ctx context.Context
* @param session sqlx.Session
* @param data *FriendRequests
* @return error
*/
func (m *customFriendRequestsModel) UpdateWithSession(ctx context.Context, session sqlx.Session, data *FriendRequests) error {
	query := fmt.Sprintf("update %s set %s where id = $1", m.table, friendRequestsRowsWithPlaceHolder)
	_, err := session.ExecCtx(ctx, query, data.Id, data.UserId, data.ReqUid, data.ReqMsg, data.ReqTime, data.HandleResult, data.HandleMsg, data.HandledAt)
	if err != nil {
		return err
	}
	return nil
}