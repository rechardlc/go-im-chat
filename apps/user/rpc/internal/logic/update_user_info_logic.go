package logic

import (
	"context"
	"errors"

	"github.com/rechardlc/go-im-chat/apps/user/models"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/internal/convert"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/internal/svc"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/user/github.com/rechardlc/go-im-chat/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrUserRequired = errors.New("user is required")
	ErrNicknameRequired = errors.New("nickname is required")
	ErrAvatarRequired = errors.New("avatar is required")
	ErrSexRequired = errors.New("sex is required")	
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
