package logic

import (
	"context"
	"errors"

	"github.com/rechardlc/go-im-chat/apps/user/models"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/internal/svc"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/user/github.com/rechardlc/go-im-chat/apps/user/rpc/user"

	pkgErr "github.com/pkg/errors"
	"github.com/rechardlc/go-im-chat/pkg/xeer"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrIdRequired   = xeer.NewError(xeer.REQUEST_PARAM_ERROR, "id is required")
	ErrTokenInvalid = xeer.NewError(xeer.REQUEST_PARAM_ERROR, "token is invalid")
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
	id := in.Id
	// 校验ID
	if id == "" {
		return nil, pkgErr.WithStack(ErrIdRequired)
	}
	// 校验Token是否有效
	// token := l.ctx.Value(ctxdata.CtxKeyJwtToken).(string)
	// if token == "" {
	// 	return nil, ErrTokenInvalid
	// }
	u, err := l.svcCtx.UsersModel.FindOne(l.ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, pkgErr.WithStack(ErrUserNotFound)
		}
		return nil, pkgErr.Wrapf(xeer.NewDBError(), "find user by id failed: %v, req: %v", err, id)
	}
	return &user.GetUserInfoResp{
		User: &user.UserEntity{
			Id:       u.Id,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
			Sex:      user.Sex(u.Sex),
		},
	}, nil
}
