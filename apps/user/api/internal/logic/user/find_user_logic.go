// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"github.com/rechardlc/go-im-chat/apps/user/api/internal/svc"
	"github.com/rechardlc/go-im-chat/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/user_client"
)

type FindUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查找用户
func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FindUserLogic) FindUser(req *types.FindUserReq) (resp *types.FindUserResp, err error) {
	// todo: add your logic here and delete this line
	var userResp *user_client.FindUserResp
	userResp, err = l.svcCtx.UserRpc.FindUser(l.ctx, &user_client.FindUserReq{
		Phone: req.Phone,
		Ids: req.Ids,
		Name: req.Name,
	})
	if err != nil {
		return nil, err
	}
	var users []types.UserEntity
	for _, user := range userResp.Users {
		users = append(users, types.UserEntity{
			Id: user.Id,
			Nickname: user.Nickname,
			Avatar: user.Avatar,
		})
	}
	return &types.FindUserResp{
		Users: users,
	}, nil
}
