package handler

import (
	"liteChat/apps/im/ws/internal/handler/conversation"
	"liteChat/apps/im/ws/internal/handler/push"
	"liteChat/apps/im/ws/internal/handler/user"
	"liteChat/apps/im/ws/internal/svc"
	"liteChat/apps/im/ws/websocket"
)

func RegisterHandlers(srv *websocket.Server, svc *svc.ServiceContext) {
	srv.AddRoutes([]websocket.Route{
		{
			Method:  "user.Online",
			Handler: user.Online(svc),
		},
		{
			Method:  "conversation.Chat",
			Handler: conversation.Chat(svc),
		},
		{
			Method:  "push",
			Handler: push.Push(svc),
		},
		{
			Method:  "conversation.MarkRead",
			Handler: conversation.MarkRead(svc),
		},
	})
}
