package models

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UsersModel = (*customUsersModel)(nil)

type (
	// UsersModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUsersModel.
	UsersModel interface {
		usersModel
		// 按昵称 / 手机号 / ids 可选条件查询用户（传了才拼进 WHERE，条件之间为 OR）。
		FindUser(ctx context.Context, name string, phone string, ids []string) ([]*Users, error)
	}

	customUsersModel struct {
		*defaultUsersModel
	}
)

// NewUsersModel returns a model for the database table.
func NewUsersModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UsersModel {
	return &customUsersModel{
		defaultUsersModel: newUsersModel(conn, c, opts...),
	}
}

// FindUser 按昵称 / 手机号 / ids 可选条件查询用户（传了才拼进 WHERE，条件之间为 OR）。
func (m *customUsersModel) FindUser(ctx context.Context, name string, phone string, ids []string) ([]*Users, error) {
	var (
		conditions []string
		args       []any
		argIdx     = 1
	)

	if name != "" {
		conditions = append(conditions, fmt.Sprintf("nickname = $%d", argIdx))
		args = append(args, name)
		argIdx++
	}
	if phone != "" {
		conditions = append(conditions, fmt.Sprintf("phone = $%d", argIdx))
		args = append(args, phone)
		argIdx++
	}
	if len(ids) > 0 {
		placeholders := make([]string, 0, len(ids))
		for _, id := range ids {
			placeholders = append(placeholders, fmt.Sprintf("$%d", argIdx))
			args = append(args, id)
			argIdx++
		}
		conditions = append(conditions, fmt.Sprintf("id IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(conditions) == 0 {
		return []*Users{}, nil
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s", usersRows, m.table, strings.Join(conditions, " OR "))

	var users []*Users
	// QueryRowsNoCacheCtx 不使用缓存，直接查询数据库。
	err := m.QueryRowsNoCacheCtx(ctx, &users, query, args...)
	if err != nil {
		return nil, err
	}
	return users, nil
}
