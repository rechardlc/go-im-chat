package logic

import (
	"context"
	"errors"

	"github.com/rechardlc/go-im-chat/apps/user/models"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/internal/convert"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/internal/svc"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/user/github.com/rechardlc/go-im-chat/apps/user/rpc/user"

	pkgErr "github.com/pkg/errors"
	"github.com/rechardlc/go-im-chat/pkg/xeer"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrUserRequired     = xeer.NewError(xeer.REQUEST_PARAM_ERROR, "用户不能为空")
	ErrNicknameRequired = xeer.NewError(xeer.REQUEST_PARAM_ERROR, "昵称不能为空")
	ErrAvatarRequired   = xeer.NewError(xeer.REQUEST_PARAM_ERROR, "头像不能为空")
	ErrSexRequired      = xeer.NewError(xeer.REQUEST_PARAM_ERROR, "性别不能为空")
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(in *user.UpdateUserInfoReq) (*user.UpdateUserInfoResp, error) {
	id := in.Id
	nickname := in.Nickname
	avatar := in.Avatar
	sex := in.Sex

	if id == "" {
		return nil, pkgErr.WithStack(ErrUserRequired)
	}
	if nickname == "" {
		return nil, pkgErr.WithStack(ErrNicknameRequired)
	}
	if avatar == "" {
		return nil, pkgErr.WithStack(ErrAvatarRequired)
	}
	if sex == 0 {
		return nil, pkgErr.WithStack(ErrSexRequired)
	}

	u, err := l.svcCtx.UsersModel.FindOne(l.ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, pkgErr.WithStack(ErrUserNotFound)
		}
		return nil, pkgErr.Wrapf(xeer.NewDBError(), "find user by id failed: %v, req: %v", err, id)
	}

	u.Nickname = nickname
	u.Avatar = avatar
	u.Sex = convert.ProtoSexToDB(sex)

	err = l.svcCtx.UsersModel.Update(l.ctx, u)
	if err != nil {
		return nil, pkgErr.Wrapf(xeer.NewDBError(), "update user failed: %v, req: %v", err, id)
	}
	return &user.UpdateUserInfoResp{
		User: &user.UserEntity{
			Id:       u.Id,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
			Phone:    u.Phone,
			Sex:      convert.DBSexToProto(u.Sex),
			Status:   convert.DBStatusToProto(u.Status),
		},
	}, nil
}
