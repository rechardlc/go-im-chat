// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package friend

import (
	"context"

	"github.com/rechardlc/go-im-chat/apps/social/api/internal/svc"
	"github.com/rechardlc/go-im-chat/apps/social/api/internal/types"

	"github.com/rechardlc/go-im-chat/apps/social/rpc/social_client"
	"github.com/rechardlc/go-im-chat/pkg/ctxdata"
	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友申请
func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInLogic) FriendPutIn(req *types.FriendPutInReq) (resp *types.FriendPutInResp, err error) {
	userID, err := ctxdata.GetJwtUserID(l.ctx)
	if err != nil {
		return resp, err
	}
	_, err = l.svcCtx.SocialRpc.FriendPutIn(l.ctx, &social_client.FriendPutInReq{
		UserId: userID,
		ReqUid: req.UserId,
		ReqMsg: req.ReqMsg,
		ReqTime: req.ReqTime,
	})
	if err != nil {
		return resp, err
	}
	return nil, nil
}
