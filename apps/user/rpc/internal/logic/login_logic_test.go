package logic

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/rechardlc/go-im-chat/apps/user/models"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/internal/svc"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/user/github.com/rechardlc/go-im-chat/apps/user/rpc/user"
	"github.com/rechardlc/go-im-chat/pkg/ctxdata"
	"github.com/rechardlc/go-im-chat/pkg/encrypt"
	"github.com/rechardlc/go-im-chat/pkg/xeer"
	pkgErr "github.com/pkg/errors"
	zeroxErrors "github.com/zeromicro/x/errors"
)

/**
 * stubUsersModel 仅用于 Login 单测，按需 stub FindOneByPhone。
 */
type stubUsersModel struct {
	findOneByPhone func(ctx context.Context, phone string) (*models.Users, error)
}

func (s *stubUsersModel) Insert(context.Context, *models.Users) (sql.Result, error) {
	panic("unexpected call: Insert")
}

func (s *stubUsersModel) FindOne(context.Context, string) (*models.Users, error) {
	panic("unexpected call: FindOne")
}

func (s *stubUsersModel) FindOneByPhone(ctx context.Context, phone string) (*models.Users, error) {
	return s.findOneByPhone(ctx, phone)
}

func (s *stubUsersModel) Update(context.Context, *models.Users) error {
	panic("unexpected call: Update")
}

func (s *stubUsersModel) Delete(context.Context, string) error {
	panic("unexpected call: Delete")
}

func (s *stubUsersModel) FindUser(context.Context, string, string, []string) ([]*models.Users, error) {
	panic("unexpected call: FindUser")
}

/**
 * newTestLoginLogic 构造带 stub 依赖的 LoginLogic。
 */
func newTestLoginLogic(findOneByPhone func(ctx context.Context, phone string) (*models.Users, error)) *LoginLogic {
	return NewLoginLogic(context.Background(), &svc.ServiceContext{
		Token: ctxdata.MustNewTokenManager(ctxdata.JWTConfig{
			Secret:             "test-secret",
			AccessTokenExpire:  time.Hour,
			RefreshTokenExpire: 24 * time.Hour,
		}),
		UsersModel: &stubUsersModel{findOneByPhone: findOneByPhone},
	})
}

func TestLoginLogic_Login(t *testing.T) {
	hashedPassword, err := encrypt.HashPassword("123456")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}

	testUser := &models.Users{
		Id:       "user-1",
		Phone:    "13800138000",
		Password: hashedPassword,
		Nickname: "tester",
		Avatar:   "https://example.com/avatar.png",
		Status:   0,
		Sex:      1,
	}

	tests := []struct {
		name       string
		req        *user.LoginReq
		findOne    func(ctx context.Context, phone string) (*models.Users, error)
		wantErr    error
		wantErrMsg string
		checkResp  func(t *testing.T, resp *user.LoginResp)
	}{
		{
			name:    "invalid phone",
			req:     &user.LoginReq{Phone: "123", Password: "123456"},
			wantErr: ErrInvalidPhone,
		},
		{
			name: "user not found",
			req:  &user.LoginReq{Phone: "13800138000", Password: "123456"},
			findOne: func(context.Context, string) (*models.Users, error) {
				return nil, models.ErrNotFound
			},
			wantErr: ErrUserNotFound,
		},
		{
			name: "database error",
			req:  &user.LoginReq{Phone: "13800138000", Password: "123456"},
			findOne: func(context.Context, string) (*models.Users, error) {
				return nil, errors.New("database unavailable")
			},
			wantErrMsg: "find user by phone failed: database unavailable",
		},
		{
			name: "wrong password",
			req:  &user.LoginReq{Phone: "13800138000", Password: "wrong-password"},
			findOne: func(context.Context, string) (*models.Users, error) {
				return testUser, nil
			},
			wantErr: ErrInvalidPassword,
		},
		{
			name: "success",
			req:  &user.LoginReq{Phone: "13800138000", Password: "123456"},
			findOne: func(_ context.Context, phone string) (*models.Users, error) {
				if phone != "13800138000" {
					t.Fatalf("FindOneByPhone phone = %q, want %q", phone, "13800138000")
				}
				return testUser, nil
			},
			checkResp: func(t *testing.T, resp *user.LoginResp) {
				if resp.Token == "" {
					t.Fatal("expected non-empty token")
				}
				if resp.User == nil {
					t.Fatal("expected user in response")
				}
				if resp.User.Id != testUser.Id {
					t.Fatalf("user id = %q, want %q", resp.User.Id, testUser.Id)
				}
				if resp.User.Phone != testUser.Phone {
					t.Fatalf("user phone = %q, want %q", resp.User.Phone, testUser.Phone)
				}
				if resp.User.Nickname != testUser.Nickname {
					t.Fatalf("user nickname = %q, want %q", resp.User.Nickname, testUser.Nickname)
				}
				if resp.User.Avatar != testUser.Avatar {
					t.Fatalf("user avatar = %q, want %q", resp.User.Avatar, testUser.Avatar)
				}
				if resp.User.Status != user.UserStatus_USER_STATUS_ACTIVE {
					t.Fatalf("user status = %v, want ACTIVE", resp.User.Status)
				}
				if resp.User.Sex != user.Sex_SEX_MALE {
					t.Fatalf("user sex = %v, want MALE", resp.User.Sex)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logic := newTestLoginLogic(tt.findOne)
			resp, err := logic.Login(tt.req)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("Login() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if tt.wantErrMsg != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Fatalf("Login() error = %v, want message containing %q", err, tt.wantErrMsg)
				}
				causeErr := pkgErr.Cause(err)
				if e, ok := causeErr.(*zeroxErrors.CodeMsg); !ok || e.Code != xeer.DB_ERROR {
					t.Fatalf("Login() cause = %v, want DB_ERROR", causeErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Login() unexpected error = %v", err)
			}
			if tt.checkResp != nil {
				tt.checkResp(t, resp)
			}
		})
	}
}
