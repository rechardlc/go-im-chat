package logic

import (
	"context"
	"errors"

	"github.com/rechardlc/go-im-chat/apps/user/models"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/internal/svc"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/user/github.com/rechardlc/go-im-chat/apps/user/rpc/user"

	"github.com/rechardlc/go-im-chat/pkg/encrypt"
	"github.com/rechardlc/go-im-chat/pkg/utils"
	"github.com/zeromicro/go-zero/core/logx"

	convert "github.com/rechardlc/go-im-chat/apps/user/rpc/internal/convert"
)


type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// 验证手机号
	if err := utils.ValidatePhone(in.Phone); err != nil {
		return nil, ErrInvalidPhone
	}

	// 查询用户
	u, err := l.svcCtx.FindOneByPhone(l.ctx, in.Phone)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		l.Errorf("find user by phone failed: %v", err)
		return nil, err
	}
	// 用户不存在
	if u == nil {
		return nil, ErrUserNotFound
	}

	// 验证密码
	if err := encrypt.ValidatePassword(u.Password, in.Password); err != nil {
		return nil, ErrInvalidPassword
	}

	// 颁发token
	token, err := l.svcCtx.Token.GenerateToken(u.Id, u.Phone)
	if err != nil {
		return nil, err
	}
	// 返回结果
	return &user.LoginResp{
		Token: token,
		User: &user.UserEntity{
			Id:       u.Id,
			Phone:    u.Phone,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
			Status:   convert.DBStatusToProto(u.Status),
			Sex:      convert.DBSexToProto(u.Sex),
		},
	}, nil
}
