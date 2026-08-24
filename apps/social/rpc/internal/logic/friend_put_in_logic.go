package logic

import (
	"context"
	"database/sql"
	"errors"
	"time"

	pkgErr "github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/rechardlc/go-im-chat/apps/social/models"
	"github.com/rechardlc/go-im-chat/apps/social/rpc/internal/svc"
	"github.com/rechardlc/go-im-chat/apps/social/rpc/social/github.com/rechardlc/go-im-chat/apps/social/rpc/social"
	"github.com/rechardlc/go-im-chat/pkg/constants"
	"github.com/rechardlc/go-im-chat/pkg/xeer"
)

var (
	ErrFriendAlreadyExists        = xeer.NewError(xeer.REQUEST_PARAM_ERROR, "已经存在好友关系")
	ErrFriendRequestAlreadyExists = xeer.NewError(xeer.REQUEST_PARAM_ERROR, "已经存在好友申请")
)

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
	// 1. 查询是否已经存在好友关系
	friend, err := l.svcCtx.FriendsModel.FindOneByUserIdAndFriendUid(l.ctx, in.UserId, in.ReqUid)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return nil, pkgErr.Wrapf(xeer.NewDBError(), "find friend failed: %v, req: userId=%v friendUid=%v", err, in.UserId, in.ReqUid)
	}
	if friend != nil {
		return nil, pkgErr.WithStack(ErrFriendAlreadyExists)
	}

	// 2. 查询是否已经存在好友申请
	friendRequest, err := l.svcCtx.FriendRequestsModel.FindOneByUserIdAndReqUid(l.ctx, in.UserId, in.ReqUid)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return nil, pkgErr.Wrapf(xeer.NewDBError(), "find friend request failed: %v, req: userId=%v reqUid=%v", err, in.UserId, in.ReqUid)
	}
	if friendRequest != nil {
		return nil, pkgErr.WithStack(ErrFriendRequestAlreadyExists)
	}

	// 3. 创建申请记录
	_, err = l.svcCtx.FriendRequestsModel.Insert(l.ctx, &models.FriendRequests{
		UserId:  in.UserId,
		ReqUid:  in.ReqUid,
		ReqMsg:  in.ReqMsg,
		ReqTime: time.Now(),
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandlerResult),
			Valid: true,
		},
	})
	if err != nil {
		return nil, pkgErr.Wrapf(xeer.NewDBError(), "insert friend request failed: %v, req: userId=%v reqUid=%v", err, in.UserId, in.ReqUid)
	}

	return &social.FriendPutInResp{}, nil
}
