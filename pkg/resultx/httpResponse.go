package resultx

import (
	"context"
	"net/http"

	"github.com/rechardlc/go-im-chat/pkg/xeer"
	pkgErrors "github.com/pkg/errors"
	zeroxErrors "github.com/zeromicro/x/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)
type Response struct { 
	Code int `json:"code"`
	Msg string `json:"message"`
	Data interface{} `json:"data"`
}

func Success(data interface{}) *Response {
	return &Response{
		Code: 0,
		Msg: "success",
		Data: data,
	}
}

func Fail(code int, err string) *Response {
	return &Response{
		Code: code,
		Msg: err,
		Data: nil,
	}
}

// OkHandler 统一成功响应格式，签名需与 httpx.SetOkHandler 一致。
func OkHandler(_ context.Context, data any) any {
	return Success(data)
}

// ErrHandler 统一错误响应格式，签名需与 httpx.SetErrorHandlerCtx 一致。
func ErrHandler(name string) func(ctx context.Context, err error) (int, any) {
	return func(ctx context.Context, err error) (int, any) {
		errcode := xeer.SERVER_COMMON_ERROR
		errmsg := xeer.Errmsg(errcode)
		caseError := pkgErrors.Cause(err)
		if e, ok := caseError.(*zeroxErrors.CodeMsg); ok {
			errcode = e.Code
			errmsg = e.Msg
		} else {
			if gstatus, ok := status.FromError(err); ok {
				errcode = int(gstatus.Code())
				errmsg = gstatus.Message()
			}
		}
		logx.WithContext(ctx).Errorf("errHandler[%s] error: %+v", name, err)
		return http.StatusBadRequest, Fail(errcode, errmsg)
	}
}