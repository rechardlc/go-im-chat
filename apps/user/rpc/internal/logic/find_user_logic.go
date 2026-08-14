package logic

import (
	"context"

	"github.com/rechardlc/go-im-chat/apps/user/rpc/internal/convert"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/internal/svc"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/user/github.com/rechardlc/go-im-chat/apps/user/rpc/user"

	pkgErr "github.com/pkg/errors"
	"github.com/rechardlc/go-im-chat/pkg/xeer"
	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserLogic) FindUser(in *user.FindUserReq) (*user.FindUserResp, error) {
	name := in.Name
	phone := in.Phone
	ids := in.Ids

	// 查询用户
	users, err := l.svcCtx.UsersModel.FindUser(l.ctx, name, phone, ids)
	if err != nil {
		return nil, pkgErr.Wrapf(xeer.NewDBError(), "find user failed: %v, req: name=%v phone=%v ids=%v", err, name, phone, ids)
	}
	// 将用户信息转换为 user.UserEntity 类型
	userEntities := make([]*user.UserEntity, 0, len(users))
	for _, u := range users {
		userEntities = append(userEntities, &user.UserEntity{
			Id:       u.Id,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
			Sex:      convert.DBSexToProto(u.Sex),
		})
	}
	return &user.FindUserResp{
		Users: userEntities,
	}, nil
}
