package utils

import (
	"errors"
	"regexp"
)

func ValidatePhone(phone string) error {
	// 直接正则匹配：长度、空、规则，一个正则一起处理
	reg := regexp.MustCompile(`^1[3-9]\d{9}$|^$`)
	if !reg.MatchString(phone) {
		return errors.New("phone is invalid")
	}
	return nil
}