package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"liteChat/apps/task/mq/internal/config"
	"liteChat/apps/task/mq/internal/handler"
	"liteChat/apps/task/mq/internal/svc"
)

var configFile = flag.String("f", "etc/task.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if err := c.SetUp(); err != nil {
		panic(err)
	}
	ctx := svc.NewServiceContext(c)
	listen := handler.NewListen(ctx)
	serviceGroup := service.NewServiceGroup()
	for _, s := range listen.Services() {
		serviceGroup.Add(s)
	}
	fmt.Println("start mq at ", c.ListenOn, "......")
	serviceGroup.Start()
	defer serviceGroup.Stop()
	fmt.Println("stop mq at ", c.ListenOn, "......")
}
