package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"liteChat/apps/im/immodels"
	"liteChat/apps/im/ws/websocket"
	"liteChat/apps/social/rpc/socialclient"
	"liteChat/apps/task/mq/internal/config"
	"liteChat/pkg/constants"
	"net/http"
)

type ServiceContext struct {
	config.Config
	WsClient websocket.Client
	*redis.Redis
	socialclient.Social
	immodels.ChatLogModel
	immodels.ConversationModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	svc := &ServiceContext{
		Config:            c,
		Redis:             redis.MustNewRedis(c.Redisx),
		Social:            socialclient.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
		ChatLogModel:      immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.DB),
		ConversationModel: immodels.MustConversationModel(c.Mongo.Url, c.Mongo.DB),
	}

	token, err := svc.GetSystemToken()
	if err != nil {
		panic(err)
	}
	header := http.Header{}
	header.Set("Authorization", token)
	svc.WsClient = websocket.NewClient(c.Ws.Host, websocket.WithClientHeader(header))

	return svc
}

func (svc *ServiceContext) GetSystemToken() (string, error) {
	return svc.Redis.Get(constants.REDIS_SYSTEM_ROOT_TOKEN)
}
