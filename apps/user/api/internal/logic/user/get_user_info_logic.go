// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/rechardlc/go-im-chat/apps/user/api/internal/svc"
	"github.com/rechardlc/go-im-chat/apps/user/api/internal/types"
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
	userID, ok := l.ctx.Value("user_id").(string)
	fmt.Println("userID", userID)
	fmt.Println("ok", ok)
	if !ok || userID == "" {
		return nil, errors.New("unauthorized")
	}
	userResp, err := l.svcCtx.UserRpc.GetUserInfo(l.ctx, &user_client.GetUserInfoReq{
		Id: userID,
	})
	if err != nil {
		return resp, err
	}
	fmt.Println("userResp", userResp)
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
