package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"liteChat/apps/user/api/internal/config"
	"liteChat/apps/user/rpc/userclient"
)

var retryPolicy = `{
	"methodConfig":[{
		"name":[{
				"service":"user.User"
			}],
		"waitForReady":true,
		"retryPolicy":{
			"maxAttempts":5,
			"initialBackoff":"0.001s",
			"maxBackoff":"0.002s",
			"backoffMultiplier":1.0,
			"retryableStatusCode":["UNKNOWN","DEADLINE_EXCEED"]
		}
	}]
}`

type ServiceContext struct {
	Config config.Config
	*redis.Redis
	userclient.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Redis:  redis.MustNewRedis(c.Redisx),
		User:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc, zrpc.WithDialOption(grpc.WithDefaultServiceConfig(retryPolicy)))),
	}
}
