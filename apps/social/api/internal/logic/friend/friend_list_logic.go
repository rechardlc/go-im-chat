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
	"github.com/rechardlc/go-im-chat/apps/user/rpc/user_client"
)

type FriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友列表
func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendListLogic) FriendList(req *types.FriendListReq) (resp *types.FriendListResp, err error) {
	// todo: add your logic here and delete this line
	userID, err := ctxdata.GetJwtUserID(l.ctx)
	if err != nil {
		return resp, err
	}
	friendListResp, err := l.svcCtx.SocialRpc.FriendList(l.ctx, &social_client.FriendListReq{
		UserId: userID,
	})
	if err != nil {
		return resp, err
	}

	if len(friendListResp.List) == 0 {
		return &types.FriendListResp{}, nil
	}
	// 提取ids
	ids := make([]string, 0, len(friendListResp.List))
	for _, v := range friendListResp.List {
		ids = append(ids, v.FriendUid)
	}

	userResp, err := l.svcCtx.UserRpc.FindUser(l.ctx, &user_client.FindUserReq{
		Ids: ids,
	})
	if err != nil {
		return resp, err
	}

	users := make(map[string]*user_client.UserEntity, len(userResp.Users))
	for _, v := range userResp.Users {
		users[v.Id] = v
	}

	friendList := make([]*types.Friends, 0, len(friendListResp.List))

	for _, v := range friendListResp.List {
		friend := &types.Friends{
			Id:        v.Id,
			FriendUid: v.FriendUid,
			Remark:    v.Remark,
		}
		if user, ok := users[v.FriendUid]; ok {
			friend.Nickname = user.Nickname
			friend.Avatar = user.Avatar
		}
		friendList = append(friendList, friend)
	}
	return &types.FriendListResp{
		List: friendList,
	}, nil
}
