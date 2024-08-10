package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/service"
	"liteChat/apps/task/mq/internal/config"
	"liteChat/apps/task/mq/internal/handler"
	"liteChat/apps/task/mq/internal/svc"
	"liteChat/pkg/configserver"
	"sync"
)

var configFile = flag.String("f", "etc/task.yaml", "the config file")
var wg sync.WaitGroup

func main() {
	flag.Parse()

	var c config.Config
	err := configserver.NewConfigServer(*configFile, configserver.NewSail(&configserver.Config{
		ETCDEndpoints:  "127.0.0.1:2379",
		ProjectKey:     "",
		Namespace:      "task",
		Configs:        "task-mq.yaml",
		ConfigFilePath: "./conf",
		LogLevel:       "DEBUG",
	})).MustLoad(&c, func(bytes []byte) error {
		var c config.Config
		_ = configserver.LoadFromJsonBytes(bytes, &c)

		wg.Add(1)
		go func(c config.Config) {
			defer wg.Done()

			Run(c)
		}(c)
		return nil
	})
	if err != nil {
		panic(err)
	}

	wg.Add(1)
	go func(c config.Config) {
		defer wg.Done()

		Run(c)
	}(c)

	wg.Wait()
}

func Run(c config.Config) {
	if err := c.SetUp(); err != nil {
		panic(err)
	}
	ctx := svc.NewServiceContext(c)
	listen := handler.NewListen(ctx)

	serviceGroup := service.NewServiceGroup()
	for _, s := range listen.Services() {
		serviceGroup.Add(s)
	}
	fmt.Println("Starting mq at ...")
	serviceGroup.Start()
}
