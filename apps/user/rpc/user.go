package main

import (
	"flag"
	"fmt"

	"github.com/rechardlc/go-im-chat/apps/user/rpc/internal/config"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/internal/server"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/internal/svc"
	"github.com/rechardlc/go-im-chat/apps/user/rpc/user/github.com/rechardlc/go-im-chat/apps/user/rpc/user"
	"github.com/rechardlc/go-im-chat/pkg/wuid"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// flag包的作用：定义参数 -> 解析参数 -> 使用参数
// -f 参数名，etc/dev/user.yaml 参数值，the config file 参数描述
var configFile = flag.String("f", "etc/dev/user.yaml", "the config file")

func main() {
	// 解析参数
	flag.Parse()

	// 加载配置
	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	if err := wuid.Init(c.Postgres.DataSource); err != nil {
		panic(err)
	}
	defer wuid.Close()

	// 创建服务上下文
	ctx := svc.NewServiceContext(c)

	// 创建RPC服务器
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServer(grpcServer, server.NewUserServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
