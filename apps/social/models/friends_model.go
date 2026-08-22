package models

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FriendsModel = (*customFriendsModel)(nil)

type (
	// FriendsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFriendsModel.
	FriendsModel interface {
		friendsModel
		FindOneByUserIdAndFriendUid(ctx context.Context, userId, friendUid string) (*Friends, error)
	}

	customFriendsModel struct {
		*defaultFriendsModel
	}
)

// NewFriendsModel returns a model for the database table.
func NewFriendsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FriendsModel {
	return &customFriendsModel{
		defaultFriendsModel: newFriendsModel(conn, c, opts...),
	}
}

func (m *customFriendsModel) FindOneByUserIdAndFriendUid(ctx context.Context, userId, friendUid string) (*Friends, error) {
	query := fmt.Sprintf("select %s from %s where user_id = $1 and friend_uid = $2 limit 1", friendsRows, m.table)
	var resp Friends
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, userId, friendUid)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
