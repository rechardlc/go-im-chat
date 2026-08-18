package logic

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"

	pkgErr "github.com/pkg/errors"
	"github.com/rechardlc/go-im-chat/apps/user/models"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/internal/svc"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/user/github.com/rechardlc/go-im-chat/apps/user/rpc/user"
	"github.com/rechardlc/go-im-chat/pkg/xeer"
	zeroxErrors "github.com/zeromicro/x/errors"
)

/**
 * stubUpdateUsersModel 仅用于 UpdateUserInfo 单测，按需 stub FindOne / Update。
 */
type stubUpdateUsersModel struct {
	findOne func(ctx context.Context, id string) (*models.Users, error)
	update  func(ctx context.Context, data *models.Users) error
}

func (s *stubUpdateUsersModel) Insert(context.Context, *models.Users) (sql.Result, error) {
	panic("unexpected call: Insert")
}

func (s *stubUpdateUsersModel) FindOne(ctx context.Context, id string) (*models.Users, error) {
	if s.findOne == nil {
		panic("unexpected call: FindOne")
	}
	return s.findOne(ctx, id)
}

func (s *stubUpdateUsersModel) FindOneByPhone(context.Context, string) (*models.Users, error) {
	panic("unexpected call: FindOneByPhone")
}

func (s *stubUpdateUsersModel) Update(ctx context.Context, data *models.Users) error {
	if s.update == nil {
		panic("unexpected call: Update")
	}
	return s.update(ctx, data)
}

func (s *stubUpdateUsersModel) Delete(context.Context, string) error {
	panic("unexpected call: Delete")
}

func (s *stubUpdateUsersModel) FindUser(context.Context, string, string, []string) ([]*models.Users, error) {
	panic("unexpected call: FindUser")
}

/**
 * newTestUpdateUserInfoLogic 构造带 stub 依赖的 UpdateUserInfoLogic。
 */
func newTestUpdateUserInfoLogic(model *stubUpdateUsersModel) *UpdateUserInfoLogic {
	return NewUpdateUserInfoLogic(context.Background(), &svc.ServiceContext{
		UsersModel: model,
	})
}

func TestUpdateUserInfoLogic_UpdateUserInfo(t *testing.T) {
	existing := &models.Users{
		Id:       "user-1",
		Phone:    "13800138000",
		Password: "hashed-password",
		Nickname: "old-nick",
		Avatar:   "https://example.com/old.png",
		Status:   0,
		Sex:      1,
	}

	tests := []struct {
		name       string
		req        *user.UpdateUserInfoReq
		model      *stubUpdateUsersModel
		wantErr    error
		wantErrMsg string
		checkResp  func(t *testing.T, resp *user.UpdateUserInfoResp)
	}{
		{
			name:    "empty id",
			req:     &user.UpdateUserInfoReq{Nickname: "n", Avatar: "a", Sex: user.Sex_SEX_MALE},
			wantErr: ErrUserRequired,
		},
		{
			name:    "empty nickname",
			req:     &user.UpdateUserInfoReq{Id: "user-1", Avatar: "a", Sex: user.Sex_SEX_MALE},
			wantErr: ErrNicknameRequired,
		},
		{
			name:    "empty avatar",
			req:     &user.UpdateUserInfoReq{Id: "user-1", Nickname: "n", Sex: user.Sex_SEX_MALE},
			wantErr: ErrAvatarRequired,
		},
		{
			name:    "unspecified sex",
			req:     &user.UpdateUserInfoReq{Id: "user-1", Nickname: "n", Avatar: "a"},
			wantErr: ErrSexRequired,
		},
		{
			name: "user not found",
			req: &user.UpdateUserInfoReq{
				Id:       "missing",
				Nickname: "n",
				Avatar:   "a",
				Sex:      user.Sex_SEX_MALE,
			},
			model: &stubUpdateUsersModel{
				findOne: func(context.Context, string) (*models.Users, error) {
					return nil, models.ErrNotFound
				},
			},
			wantErr: ErrUserNotFound,
		},
		{
			name: "find one database error",
			req: &user.UpdateUserInfoReq{
				Id:       "user-1",
				Nickname: "n",
				Avatar:   "a",
				Sex:      user.Sex_SEX_MALE,
			},
			model: &stubUpdateUsersModel{
				findOne: func(context.Context, string) (*models.Users, error) {
					return nil, errors.New("database unavailable")
				},
			},
			wantErrMsg: "find user by id failed: database unavailable",
		},
		{
			name: "preserves phone password status",
			req: &user.UpdateUserInfoReq{
				Id:       "user-1",
				Nickname: "new-nick",
				Avatar:   "https://example.com/new.png",
				Sex:      user.Sex_SEX_FEMALE,
			},
			model: &stubUpdateUsersModel{
				findOne: func(_ context.Context, id string) (*models.Users, error) {
					if id != existing.Id {
						t.Fatalf("FindOne id = %q, want %q", id, existing.Id)
					}
					copied := *existing
					return &copied, nil
				},
				update: func(_ context.Context, data *models.Users) error {
					if data.Phone != existing.Phone {
						t.Fatalf("Update phone = %q, want %q", data.Phone, existing.Phone)
					}
					if data.Password != existing.Password {
						t.Fatalf("Update password = %q, want %q", data.Password, existing.Password)
					}
					if data.Status != existing.Status {
						t.Fatalf("Update status = %d, want %d", data.Status, existing.Status)
					}
					if data.Nickname != "new-nick" {
						t.Fatalf("Update nickname = %q, want %q", data.Nickname, "new-nick")
					}
					if data.Avatar != "https://example.com/new.png" {
						t.Fatalf("Update avatar = %q, want %q", data.Avatar, "https://example.com/new.png")
					}
					if data.Sex != 2 {
						t.Fatalf("Update sex = %d, want 2", data.Sex)
					}
					return nil
				},
			},
			checkResp: func(t *testing.T, resp *user.UpdateUserInfoResp) {
				if resp.User == nil {
					t.Fatal("expected user in response")
				}
				if resp.User.Id != existing.Id {
					t.Fatalf("user id = %q, want %q", resp.User.Id, existing.Id)
				}
				if resp.User.Phone != existing.Phone {
					t.Fatalf("user phone = %q, want %q", resp.User.Phone, existing.Phone)
				}
				if resp.User.Nickname != "new-nick" {
					t.Fatalf("user nickname = %q, want %q", resp.User.Nickname, "new-nick")
				}
				if resp.User.Avatar != "https://example.com/new.png" {
					t.Fatalf("user avatar = %q, want %q", resp.User.Avatar, "https://example.com/new.png")
				}
				if resp.User.Status != user.UserStatus_USER_STATUS_ACTIVE {
					t.Fatalf("user status = %v, want ACTIVE", resp.User.Status)
				}
				if resp.User.Sex != user.Sex_SEX_FEMALE {
					t.Fatalf("user sex = %v, want FEMALE", resp.User.Sex)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logic := newTestUpdateUserInfoLogic(tt.model)
			resp, err := logic.UpdateUserInfo(tt.req)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("UpdateUserInfo() error = %v, want %v", err, tt.wantErr)
				}
				return
			}

			if tt.wantErrMsg != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Fatalf("UpdateUserInfo() error = %v, want message containing %q", err, tt.wantErrMsg)
				}
				causeErr := pkgErr.Cause(err)
				if e, ok := causeErr.(*zeroxErrors.CodeMsg); !ok || e.Code != xeer.DB_ERROR {
					t.Fatalf("UpdateUserInfo() cause = %v, want DB_ERROR", causeErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("UpdateUserInfo() unexpected error = %v", err)
			}
			if tt.checkResp != nil {
				tt.checkResp(t, resp)
			}
		})
	}
}
