package interceptor
import (
	"context"
	pkgErrors "github.com/pkg/errors"
	zeroxErrors "github.com/zeromicro/x/errors"
	gCodes "google.golang.org/grpc/codes" // 用于设置错误码
	gStatus "google.golang.org/grpc/status" // 用于设置错误信息
	"google.golang.org/grpc"
	"github.com/zeromicro/go-zero/core/logx"
)

func LogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err == nil {
		return resp, nil
	}
	logx.WithContext(ctx).Errorf("【RPC SERVER ERROR】: %v, req: %v, resp: %v", err, req, resp)
	causeErr := pkgErrors.Cause(err)
	if e, ok := causeErr.(*zeroxErrors.CodeMsg); ok {
		err = gStatus.Error(gCodes.Code(e.Code), e.Msg)
	}
	return resp, err
}