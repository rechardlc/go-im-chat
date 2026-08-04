package convert

import (
	"github.com/rechardlc/go-im-chat/apps/user/rpc/user/github.com/rechardlc/go-im-chat/apps/user/rpc/user"
)

// DB status: 0=正常, 1=禁用（见 deploy/sql/user.sql）
const (
	DbStatusActive   int64 = 0
	DbStatusDisabled int64 = 1
)

// DB sex: 0=未知, 1=男, 2=女（见 deploy/sql/user.sql）
const (
	DbSexUnknown int64 = 0
	DbSexMale    int64 = 1
	DbSexFemale  int64 = 2
)

// ProtoStatusToDB 将 proto UserStatus 转为 DB status。
func ProtoStatusToDB(s user.UserStatus) int64 {
	switch s {
	case user.UserStatus_USER_STATUS_ACTIVE:
		return DbStatusActive
	case user.UserStatus_USER_STATUS_DISABLED:
		return DbStatusDisabled
	default:
		return DbStatusActive
	}
}

// DBStatusToProto 将 DB status 转为 proto UserStatus。
func DBStatusToProto(s int64) user.UserStatus {
	switch s {
	case DbStatusActive:
		return user.UserStatus_USER_STATUS_ACTIVE
	case DbStatusDisabled:
		return user.UserStatus_USER_STATUS_DISABLED
	default:
		return user.UserStatus_USER_STATUS_UNSPECIFIED
	}
}

// ProtoSexToDB 将 proto Sex 转为 DB sex。
func ProtoSexToDB(s user.Sex) int64 {
	switch s {
	case user.Sex_SEX_MALE:
		return DbSexMale
	case user.Sex_SEX_FEMALE:
		return DbSexFemale
	default:
		// SEX_UNSPECIFIED、SEX_UNKNOWN 及未知值均视为未知
		return DbSexUnknown
	}
}

// DBSexToProto 将 DB sex 转为 proto Sex。
func DBSexToProto(s int64) user.Sex {
	switch s {
	case DbSexMale:
		return user.Sex_SEX_MALE
	case DbSexFemale:
		return user.Sex_SEX_FEMALE
	default:
		return user.Sex_SEX_UNKNOWN
	}
}
