// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"github.com/rechardlc/go-im-chat/apps/user/api/internal/svc"
	"github.com/rechardlc/go-im-chat/apps/user/api/internal/types"
	"github.com/rechardlc/go-im-chat/pkg/ctxdata"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/user_client"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户信息
func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo() (resp *types.GetUserInfoResp, err error) {
	userID, err := ctxdata.GetJwtUserID(l.ctx)
	if err != nil {
		return nil, err
	}
	userResp, err := l.svcCtx.UserRpc.GetUserInfo(l.ctx, &user_client.GetUserInfoReq{
		Id: userID,
	})
	if err != nil {
		return nil, err
	}
	return &types.GetUserInfoResp{
		User: types.UserEntity{
			Id: userResp.User.Id,
			Nickname: userResp.User.Nickname,
			Avatar: userResp.User.Avatar,
			Phone: userResp.User.Phone,
			Status: int32(userResp.User.Status),
			Sex: int32(userResp.User.Sex),
		},
	}, nil
}
