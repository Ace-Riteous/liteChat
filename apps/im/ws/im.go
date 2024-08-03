package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"liteChat/apps/im/ws/internal/config"
	"liteChat/apps/im/ws/internal/handler"
	"liteChat/apps/im/ws/internal/svc"
	"liteChat/apps/im/ws/websocket"
	"time"
)

var configFile = flag.String("f", "etc/im.yaml", "the etc file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if err := c.SetUp(); err != nil {
		panic(err)
	}
	ctx := svc.NewServiceContext(c)
	srv := websocket.NewServer(c.ListenOn,
		websocket.WithServerAuthentication(handler.NewJwtAuth(ctx)),
		websocket.WithServerPatten("/ws"),
		websocket.WithServerAck(websocket.RigorAck),
		websocket.WithServerMaxConnectionIdle(10*time.Second),
	)

	handler.RegisterHandlers(srv, ctx)

	fmt.Println("start websocket at ", c.ListenOn, "......")
	srv.Start()

	defer srv.Stop()
}
