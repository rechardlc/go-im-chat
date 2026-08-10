package logic

import (
	"context"

	"github.com/rechardlc/go-im-chat/apps/user/models"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/internal/convert"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/internal/svc"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/user/github.com/rechardlc/go-im-chat/apps/user/rpc/user"

	"github.com/rechardlc/go-im-chat/pkg/xeer"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrUserRequired = xeer.NewError(xeer.REQUEST_PARAM_ERROR, "用户不能为空")
	ErrNicknameRequired = xeer.NewError(xeer.REQUEST_PARAM_ERROR, "昵称不能为空")
	ErrAvatarRequired = xeer.NewError(xeer.REQUEST_PARAM_ERROR, "头像不能为空")
	ErrSexRequired = xeer.NewError(xeer.REQUEST_PARAM_ERROR, "性别不能为空")	
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
		return nil, ErrUserRequired
	}
	if nickname == "" {
		return nil, ErrUserRequired
	}
	if avatar == "" {
		return nil, ErrUserRequired
	}
	if sex == 0 {
		return nil, ErrUserRequired
	}

	u := &models.Users{
		Id:       id,
		Nickname: nickname,
		Avatar:   avatar,
		Sex:      convert.ProtoSexToDB(sex),
	}
	// 更新用户信息
	err := l.svcCtx.UsersModel.Update(l.ctx, u)
	if err != nil {
		return nil, err
	}
	return &user.UpdateUserInfoResp{
		User: &user.UserEntity{
			Id:       id,
			Nickname: nickname,
			Avatar:   avatar,
			Sex:      convert.DBSexToProto(u.Sex),
			Status:   convert.DBStatusToProto(u.Status),
		},
	}, nil
}
