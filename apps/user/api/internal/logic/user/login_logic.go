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

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户登录
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	// todo: add your logic here and delete this line
	loginResp, err := l.svcCtx.UserRpc.Login(l.ctx, &user_client.LoginReq{
		Phone: req.Phone,
		Password: req.Password,
	})
	if err != nil {
		return resp, err
	}
	return &types.LoginResp{
		Token: loginResp.Token,
		Expire: loginResp.Expire,
		User: types.UserEntity{
			Id: loginResp.User.Id,
			Nickname: loginResp.User.Nickname,
			Avatar: loginResp.User.Avatar,
			Phone: loginResp.User.Phone,
			Status: int32(loginResp.User.Status),
			Sex: int32(loginResp.User.Sex),
		},
	}, nil
}
