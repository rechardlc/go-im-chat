// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package friend

import (
	"context"

	"github.com/rechardlc/go-im-chat/apps/social/api/internal/svc"
	"github.com/rechardlc/go-im-chat/apps/social/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendsOnlineLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友在线
func NewFriendsOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendsOnlineLogic {
	return &FriendsOnlineLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendsOnlineLogic) FriendsOnline(req *types.FriendsOnlineReq) (resp *types.FriendOnlineResp, err error) {
	// todo: add your logic here and delete this line

	return
}
