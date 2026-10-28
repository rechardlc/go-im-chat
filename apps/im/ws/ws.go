package main

import (
	"flag"
	"fmt"

	"github.com/rechardlc/go-im-chat/apps/im/ws/internal/config"
	"github.com/rechardlc/go-im-chat/apps/im/ws/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file of im service")

func main() {
	flag.Parse()

	var c config.Config

	conf.MustLoad(*configFile, &c, conf.UseEnv())
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()
	svc.NewServiceContext(c)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
