package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"liteChat/apps/im/rpc/imclient"
	"liteChat/apps/social/api/internal/config"
	"liteChat/apps/social/rpc/socialclient"
	"liteChat/apps/user/rpc/userclient"
)

type ServiceContext struct {
	Config config.Config

	socialclient.Social
	userclient.User
	imclient.Im
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Social: socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
		User:   userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Im:     imclient.NewIm(zrpc.MustNewClient(c.ImRpc)),
	}
}
