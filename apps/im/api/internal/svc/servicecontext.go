package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"liteChat/apps/im/api/internal/config"
	"liteChat/apps/im/rpc/imclient"
)

type ServiceContext struct {
	Config config.Config
	imclient.Im
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Im:     imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
	}
}
