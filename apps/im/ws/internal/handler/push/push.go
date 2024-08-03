package push

import (
	"github.com/mitchellh/mapstructure"
	"liteChat/apps/im/ws/internal/svc"
	"liteChat/apps/im/ws/websocket"
	"liteChat/apps/im/ws/ws"
	"liteChat/pkg/constants"
)

func Push(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var data ws.Push
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			_ = srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
		switch data.ChatType {
		case constants.SingleChatType:
			_ = Single(srv, &data, data.RecvId)
		case constants.GroupChatType:
			_ = Group(srv, &data)
		}

	}
}

func Single(srv *websocket.Server, data *ws.Push, recvId string) error {
	rconn := srv.GetConn(recvId)
	if rconn == nil {
		//todo: 离线状态处理
		return nil
	}
	srv.Infof("push msg %v", data)

	return srv.Send(websocket.NewMessage(data.SendId, &ws.Chat{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendTime:       data.SendTime,
		Msg: ws.Msg{
			MType:   data.MType,
			Content: data.Content,
		},
	}), rconn)
}

func Group(srv *websocket.Server, data *ws.Push) error {
	for _, id := range data.RecvIds {
		func(id string) {
			srv.Schedule(func() {
				_ = Single(srv, data, id)
			})
		}(id)
	}
	return nil
}
