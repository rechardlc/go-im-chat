package models

import (
	"context"
	"fmt"
	"strings"

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
		InsertBatchsWithSession(ctx context.Context, session sqlx.Session, data []*Friends) error
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

func (m *customFriendsModel) InsertBatchsWithSession(ctx context.Context, session sqlx.Session, data []*Friends) error {
	if len(data) == 0 {
		return nil
	}
	const colsLen = 4
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("insert into %s (%s) values ", m.table, friendsRowsExpectAutoSet))
	values := make([]interface{}, 0, len(data)*colsLen)
	argIdx := 1
	for i, item := range data {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("($%d, $%d, $%d, $%d)", argIdx, argIdx+1, argIdx+2, argIdx+3))
		argIdx += colsLen
		values = append(values, item.UserId, item.FriendUid, item.Remark, item.AddSource)
	}
	_, err := session.ExecCtx(ctx, sb.String(), values...)
	return err
}