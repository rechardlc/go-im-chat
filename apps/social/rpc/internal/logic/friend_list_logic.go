package logic

import (
	"context"

	"github.com/rechardlc/go-im-chat/apps/social/rpc/internal/svc"
	"github.com/rechardlc/go-im-chat/apps/social/rpc/social/github.com/rechardlc/go-im-chat/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
	pkgErr "github.com/pkg/errors"
	"github.com/rechardlc/go-im-chat/pkg/xeer"
)

type FriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendListLogic) FriendList(in *social.FriendListReq) (*social.FriendListResp, error) {
	// todo: add your logic here and delete this line
	friends, err := l.svcCtx.FriendsModel.ListByUserId(l.ctx, in.UserId)
	if err != nil {
		return nil, pkgErr.Wrapf(xeer.NewError(xeer.DB_ERROR, "获取好友列表失败"), "获取好友列表失败: %s，异常：%s", in.UserId, err.Error())
	}
	list := make([]*social.Friends, len(friends))
	for i, friend := range friends {
		list[i] = &social.Friends{
			Id: int32(friend.Id),
			UserId: friend.UserId,
			Remark: friend.Remark,
			AddSource: int32(friend.AddSource.Int64),
			FriendUid: friend.FriendUid,
		}
	}
	return &social.FriendListResp{List: list}, nil
}
