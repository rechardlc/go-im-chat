package user_client

import (
	"github.com/rechardlc/go-im-chat/apps/user/rpc/user/github.com/rechardlc/go-im-chat/apps/user/rpc/user"
)

// goctl 生成的 user.go 只 alias 了 message，枚举需在此导出，供 API 层转换。
type (
	Sex        = user.Sex
	UserStatus = user.UserStatus
)
