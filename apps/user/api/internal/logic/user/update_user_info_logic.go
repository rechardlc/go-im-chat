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

type UpdateUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新用户信息
func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(req *types.UpdateUserInfoReq) (resp *types.UpdateUserInfoResp, err error) {
	// todo: add your logic here and delete this line
	userID, err := ctxdata.GetJwtUserID(l.ctx)
	if err != nil {
		return nil, err
	}
	var userResp *user_client.UpdateUserInfoResp
	userResp, err = l.svcCtx.UserRpc.UpdateUserInfo(l.ctx, &user_client.UpdateUserInfoReq{
		Id: userID,
		Nickname: req.Nickname,
		Avatar: req.Avatar,
		Sex: user_client.Sex(req.Sex),
	})
	if err != nil {
		return nil, err
	}
	return &types.UpdateUserInfoResp{
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
