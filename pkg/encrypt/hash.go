package encrypt

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const (
	// bcryptCost 哈希计算成本，值越大越安全但越慢。
	bcryptCost = bcrypt.DefaultCost
	// MaxPasswordLen bcrypt 支持的明文密码最大长度（字节）。
	MaxPasswordLen = 72
)

// ErrPasswordTooLong 明文密码超过 bcrypt 允许的最大长度。
var ErrPasswordTooLong = errors.New("password exceeds maximum length")

/**
 * HashPassword 使用 bcrypt 对明文密码进行哈希。
 * @param plain 明文密码
 * @returns 可存入数据库的 bcrypt 哈希字符串
 */
func HashPassword(plain string) (string, error) {
	if len(plain) > MaxPasswordLen {
		return "", fmt.Errorf("%w: max %d bytes", ErrPasswordTooLong, MaxPasswordLen)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcryptCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

/**
 * ValidatePassword 校验明文密码是否与 bcrypt 哈希匹配。
 * @param hashed 数据库中的 bcrypt 哈希
 * @param plain 用户输入的明文密码
 * @returns 匹配返回 nil，否则返回 ErrPasswordTooLong
 */
func ValidatePassword(hashed, plain string) error {
	if len(plain) > MaxPasswordLen {
		return fmt.Errorf("%w: max %d bytes", ErrPasswordTooLong, MaxPasswordLen)
	}

	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
