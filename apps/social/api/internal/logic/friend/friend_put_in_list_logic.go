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

type FriendPutInListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友申请列表
func NewFriendPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInListLogic {
	return &FriendPutInListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInListLogic) FriendPutInList(req *types.FriendPutInListReq) (resp *types.FriendPutInListResp, err error) {
	// todo: add your logic here and delete this line
	userID, err := ctxdata.GetJwtUserID(l.ctx)
	if err != nil {
		return resp, err
	}
	friendPutInListResp, err := l.svcCtx.SocialRpc.FriendPutInList(l.ctx, &social_client.FriendPutInListReq{
		UserId: userID,
	})
	if err != nil {
		return resp, err
	}
	var friendPutInList []*types.FriendRequests
	for _, friendPutIn := range friendPutInListResp.List {
		friendPutInList = append(friendPutInList, &types.FriendRequests{
			Id: int64(friendPutIn.Id),
			UserId: friendPutIn.UserId,
			ReqUid: friendPutIn.ReqUid,
			ReqMsg: friendPutIn.ReqMsg,
		})
	}
	return &types.FriendPutInListResp{
		List: friendPutInList,
	}, nil
}
