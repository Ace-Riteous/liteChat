package svc

import (
	"liteChat/apps/im/immodels"
	"liteChat/apps/im/rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config

	immodels.ChatLogModel
	immodels.ConversationModel
	immodels.ConversationsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:             c,
		ChatLogModel:       immodels.MustChatLogModel(c.Mongo.Url, c.Mongo.DB),
		ConversationModel:  immodels.MustConversationModel(c.Mongo.Url, c.Mongo.DB),
		ConversationsModel: immodels.MustConversationsModel(c.Mongo.Url, c.Mongo.DB),
	}
}
