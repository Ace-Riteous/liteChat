package config

import (
	"github.com/zeromicro/go-zero/core/service"
)

type Config struct {
	service.ServiceConf

	ListenOn string

	JwtAuth struct {
		AccessSecret string
		AccessExpire int64
	}
	Mongo struct {
		Url string
		DB  string
	}
	MsgChatTransfer struct {
		Topic string
		Addrs []string
	}
}
