// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"github.com/rechardlc/go-im-chat/apps/user/api/internal/svc"
	"github.com/rechardlc/go-im-chat/apps/user/api/internal/types"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/user_client"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户注册
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	// todo: add your logic here and delete this line
	registerResp, err := l.svcCtx.UserRpc.Register(l.ctx, &user_client.RegisterReq{
		Phone: req.Phone,
		Password: req.Password,
		Nickname: req.Nickname,
	})
	if err != nil {
		return resp, err
	}
	return &types.RegisterResp{
		Token: registerResp.Token,
		Expire: registerResp.Expire,
		User: types.UserEntity{
			Id: registerResp.User.Id,
			Nickname: registerResp.User.Nickname,
			Phone: registerResp.User.Phone,
			Status: int32(registerResp.User.Status),
			Sex: int32(registerResp.User.Sex),
		},
	}, nil
}
