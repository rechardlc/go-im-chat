// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package friend

import (
	"context"

	"github.com/rechardlc/go-im-chat/apps/social/api/internal/svc"
	"github.com/rechardlc/go-im-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/rechardlc/go-im-chat/pkg/ctxdata"
	"github.com/rechardlc/go-im-chat/apps/social/rpc/social_client"
)

type FriendPutInHandleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友申请处理
func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInHandleLogic) FriendPutInHandle(req *types.FriendPutInHandleReq) (resp *types.FriendPutInHandleResp, err error) {
	// todo: add your logic here and delete this line
	userID, err := ctxdata.GetJwtUserID(l.ctx)
	if err != nil {
		return resp, err
	}
	_, err = l.svcCtx.SocialRpc.FriendPutInHandle(l.ctx, &social_client.FriendPutInHandleReq{
		UserId: userID,
		FriendReqId: req.FriendReqId,
		HandleResult: req.HandleResult,
	})
	if err != nil {
		return resp, err
	}
	return &types.FriendPutInHandleResp{
	}, nil
}