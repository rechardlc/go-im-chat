package logic

import (
	"context"
	"errors"
	"strings"

	"github.com/rechardlc/go-im-chat/apps/user/models"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/internal/convert"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/internal/svc"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/user/github.com/rechardlc/go-im-chat/apps/user/rpc/user"

	"github.com/rechardlc/go-im-chat/pkg/encrypt"
	"github.com/rechardlc/go-im-chat/pkg/utils"
	"github.com/rechardlc/go-im-chat/pkg/wuid"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	// 1. 需要手机号合法性验证
	phone := strings.TrimSpace(in.Phone)
	if len(phone) == 0 {
		return nil, errors.New("phone is required")
	}
	if err := utils.ValidatePhone(phone); err != nil {
		return nil, err
	}
	// 2. 检查手机号是否已经注册
	userInfo, err := l.svcCtx.UsersModel.FindOneByPhone(l.ctx, phone)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		l.Errorf("find user by phone failed: %v", err)
		return nil, err
	}
	if userInfo != nil {
		return nil, ErrUserAlreadyExists
	}
	// 3. 创建用户
	userId, err := wuid.SafeNext()
	if err != nil {
		return nil, err
	}
	userInfo = &models.Users{
		Id:       userId,
		Phone:    phone,
		Nickname: strings.TrimSpace(in.Nickname),
		Avatar:   strings.TrimSpace(in.Avatar),
		Status:   convert.ProtoStatusToDB(user.UserStatus_USER_STATUS_ACTIVE),
		Sex:      convert.ProtoSexToDB(in.GetSex()),
	}

	// 4. 如果密码不为空，则生成加密密码
	if len(in.Password) > 0 {
		genPwd, err := encrypt.HashPassword(in.Password)
		if err != nil {
			return nil, err
		}
		userInfo.Password = genPwd
	}
	// 5. 插入用户
	_, err = l.svcCtx.UsersModel.Insert(l.ctx, userInfo)
	if err != nil {
		return nil, err
	}
	// 6. 生成token
	token, err := l.svcCtx.Token.GenerateToken(userId, phone)
	if err != nil {
		return nil, err
	}
	return &user.RegisterResp{
		Token: token,
		User: &user.UserEntity{
			Id:       userInfo.Id,
			Phone:    userInfo.Phone,
			Nickname: userInfo.Nickname,
			Avatar:   userInfo.Avatar,
			Status:   convert.DBStatusToProto(userInfo.Status),
			Sex:      convert.DBSexToProto(userInfo.Sex),
		},
	}, nil
}
