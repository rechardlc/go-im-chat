package logic

import (
	"context"
	"errors"
	"time"

	"github.com/rechardlc/go-im-chat/apps/social/models"
	"github.com/rechardlc/go-im-chat/apps/social/rpc/internal/svc"
	"github.com/rechardlc/go-im-chat/apps/social/rpc/social/github.com/rechardlc/go-im-chat/apps/social/rpc/social"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/rechardlc/go-im-chat/pkg/constants"
	"database/sql"
)

var ErrFriendAlreadyExists = errors.New("已经存在好友关系")
var ErrFriendRequestAlreadyExists = errors.New("已经存在好友申请")
var ErrUserNotFound = errors.New("用户不存在")

type FriendPutInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInLogic {
	return &FriendPutInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInLogic) FriendPutIn(in *social.FriendPutInReq) (*social.FriendPutInResp, error) {
	// todo: add your logic here and delete this line
	// 判断用户是否存在
	// user, err := l.svcCtx.UsersModel.FindOne(l.ctx, in.UserId)
	// if err != nil && !errors.Is(err, models.ErrNotFound) {
	// 	return nil, err
	// }
	// if user == nil {
	// 	return nil, ErrUserNotFound
	// }
	// 1. 查询是否已经存在好友关系
	friend, err := l.svcCtx.FriendsModel.FindOneByUserIdAndFriendUid(l.ctx, in.UserId, in.ReqUid)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return nil, err
	}
	if friend != nil {
		return nil, ErrFriendAlreadyExists
	}
	// 2. 查询是否已经存在好友申请
	friendRequest, err := l.svcCtx.FriendRequestsModel.FindOneByUserIdAndReqUid(l.ctx, in.UserId, in.ReqUid)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return nil, err
	}
	if friendRequest != nil {
		return nil, ErrFriendRequestAlreadyExists
	}
	// 3. 创建申请记录
	_, err = l.svcCtx.FriendRequestsModel.Insert(l.ctx, &models.FriendRequests{
		UserId: in.UserId,
		ReqUid: in.ReqUid,
		ReqMsg: in.ReqMsg,
		ReqTime: time.Now(),
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandlerResult),
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return &social.FriendPutInResp{}, nil
}
