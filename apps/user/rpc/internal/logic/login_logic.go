package logic

import (
	"context"
	"errors"
	"github.com/rechardlc/go-im-chat/pkg/xeer"
	"github.com/rechardlc/go-im-chat/apps/user/models"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/internal/svc"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/user/github.com/rechardlc/go-im-chat/apps/user/rpc/user"
	pkgErr "github.com/pkg/errors" // 用于包装错误
	"github.com/rechardlc/go-im-chat/pkg/encrypt"
	"github.com/rechardlc/go-im-chat/pkg/utils"
	"github.com/zeromicro/go-zero/core/logx"

	convert "github.com/rechardlc/go-im-chat/apps/user/rpc/internal/convert"
)

var (
	ErrInvalidPhone = xeer.NewError(xeer.REQUEST_PARAM_ERROR, "无效的电话号码")
	ErrUserNotFound = xeer.NewError(xeer.REQUEST_PARAM_ERROR, "用户不存在")
	ErrInvalidPassword = xeer.NewError(xeer.REQUEST_PARAM_ERROR, "无效的密码")
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
			return nil, pkgErr.WithStack(ErrUserNotFound)
		}
		return nil, pkgErr.Wrapf(xeer.NewDBError(), "find user by phone failed: %v, req: %v", err, in.Phone)
	}
	// 用户不存在
	if u == nil {
		return nil, pkgErr.WithStack(ErrUserNotFound)
	}

	// 验证密码
	if err := encrypt.ValidatePassword(u.Password, in.Password); err != nil {
		return nil, pkgErr.WithStack(ErrInvalidPassword)
	}

	// 颁发token
	token, err := l.svcCtx.Token.GenerateToken(u.Id, u.Phone)
	if err != nil {
		return nil, pkgErr.Wrapf(xeer.NewRequestParamError(), "generate token failed: %v, req: %v", err, in.Phone)
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
