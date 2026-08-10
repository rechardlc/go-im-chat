package xeer

import (
	"github.com/zeromicro/x/errors"
)

/*
* 创建错误
* @param errcode 错误码
* @param errmsg 错误信息
* @return error 错误
*/
func NewError(errcode int, errmsg string) error {
	return errors.New(errcode, errmsg)
}
/*
* 创建数据库错误
* @return error 错误
*/

func NewDBError() error {
	return NewError(DB_ERROR, Errmsg(DB_ERROR))
}
/*
* 创建请求参数错误
* @return error 错误
*/
func NewRequestParamError() error {
	return NewError(REQUEST_PARAM_ERROR, Errmsg(REQUEST_PARAM_ERROR))
}

/*
* 创建服务器异常错误
* @return error 错误
*/
func NewServerCommonError() error {
	return NewError(SERVER_COMMON_ERROR, Errmsg(SERVER_COMMON_ERROR))
}