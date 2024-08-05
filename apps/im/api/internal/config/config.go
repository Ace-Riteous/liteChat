package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	JwtAuth struct {
		AccessSecret string
	}
	ImRpc     zrpc.RpcClientConf
	UserRpc   zrpc.RpcClientConf
	SocialPrc zrpc.RpcClientConf
}
