package logic

import (
	"context"

	"github.com/rechardlc/go-im-chat/apps/social/rpc/internal/svc"
	"github.com/rechardlc/go-im-chat/apps/social/rpc/social/github.com/rechardlc/go-im-chat/apps/social/rpc/social"

	"database/sql"
	"time"

	pkgErr "github.com/pkg/errors"
	"github.com/rechardlc/go-im-chat/apps/social/models"
	"github.com/rechardlc/go-im-chat/pkg/constants"
	"github.com/rechardlc/go-im-chat/pkg/xeer"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	ErrFriendRequestAlreadyHandled  = xeer.NewError(xeer.REQUEST_PARAM_ERROR, "好友申请记录已经处理")
	ErrFriendRequestAlreadyPassed   = xeer.NewError(xeer.REQUEST_PARAM_ERROR, "好友申请记录已经通过")
	ErrFriendRequestAlreadyRefused  = xeer.NewError(xeer.REQUEST_PARAM_ERROR, "好友申请记录已经拒绝")
	ErrFriendRequestAlreadyCanceled = xeer.NewError(xeer.REQUEST_PARAM_ERROR, "好友申请记录已经取消")
	ErrFriendRequestInvalidStatus   = xeer.NewError(xeer.REQUEST_PARAM_ERROR, "好友申请记录状态不合法")
)

type FriendPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInHandleLogic {
	return &FriendPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendPutInHandleLogic) FriendPutInHandle(in *social.FriendPutInHandleReq) (*social.FriendPutInHandleResp, error) {
	// todo: add your logic here and delete this line
	// 好友申请处理
	// 1. 获取好友申请记录
	friendRequest, err := l.svcCtx.FriendRequestsModel.FindOne(l.ctx, int64(in.FriendReqId))
	if err != nil {
		return nil, pkgErr.Wrapf(xeer.NewError(xeer.DB_ERROR, "获取好友申请记录失败"), "获取好友申请记录失败: %s，异常：%s", in.FriendReqId, err.Error())
	}
	if friendRequest == nil {
		return nil, pkgErr.Wrapf(xeer.NewError(xeer.REQUEST_PARAM_ERROR, "好友申请记录不存在"), "好友申请记录不存在: %s，异常：%s", in.FriendReqId, err.Error())
	}
	// 2. 判断好友申请记录是否已经处理
	switch friendRequest.HandleResult.Int64 {
	case int64(constants.NoHandlerResult):
		break
	case int64(constants.Pass):
		return nil, ErrFriendRequestAlreadyPassed
	case int64(constants.Refuse):
		return nil, ErrFriendRequestAlreadyRefused
	case int64(constants.Cancel):
		return nil, ErrFriendRequestAlreadyCanceled
	default:
		return nil, ErrFriendRequestInvalidStatus
	}
	// 3. 处理好友申请
	friendRequest.HandleResult = sql.NullInt64{
		Int64: int64(in.HandleResult),
		Valid: true,
	}
	friendRequest.HandledAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	// 4. 更新好友申请记录: 事务更新
	// 4.1 更新好友申请记录
	// 4.2 创建好友关系 2条相互的记录
	err = l.svcCtx.FriendRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// 4.1 更新好友申请记录
		if err := l.svcCtx.FriendRequestsModel.UpdateWithSession(ctx, session, friendRequest); err != nil {
			return pkgErr.Wrapf(xeer.NewDBError(), "更新好友申请记录失败: %s，异常：%s", in.FriendReqId, err.Error())
		}
		// 4.2 创建好友关系 2条相互的记录
		friends := []*models.Friends{
			{
				UserId:    friendRequest.UserId,
				FriendUid: friendRequest.ReqUid,
			},
			{
				UserId:    friendRequest.ReqUid,
				FriendUid: friendRequest.UserId,
			},
		}
		return l.svcCtx.FriendsModel.InsertBatchsWithSession(ctx, session, friends)
	})
	if err != nil {
		return nil, pkgErr.Wrapf(xeer.NewDBError(), "处理好友申请失败: %s，异常：%s", in.FriendReqId, err.Error())
	}
	return nil, nil
}
