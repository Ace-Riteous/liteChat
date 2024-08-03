package conversation

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"liteChat/apps/im/ws/internal/svc"
	"liteChat/apps/im/ws/websocket"
	"liteChat/apps/im/ws/ws"
	"liteChat/apps/task/mq/mq"
	"liteChat/pkg/constants"
	"liteChat/pkg/wuid"
	"time"
)

func Chat(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var data ws.Chat
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			_ = srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
		if data.ConversationId == "" {
			switch data.ChatType {
			case constants.SingleChatType:
				data.ConversationId = wuid.CombineId(conn.Uid, data.RecvId)
				//c := logic.NewConversation(context.Background(), srv, svc)
				//err := c.SingleChat(&data, conn.Uid)
				//if err != nil {
				//	_ = srv.Send(websocket.NewErrMessage(err), conn)
				//	return
				//}
				//_ = srv.SendByUserId(websocket.NewMessage(conn.Uid, ws.Chat{
				//	ConversationId: data.ConversationId,
				//	ChatType:       data.ChatType,
				//	SendId:         data.SendId,
				//	RecvId:         data.RecvId,
				//	SendTime:       time.Now().UnixNano(),
				//	Msg:            data.Msg,
				//}), data.RecvId)
			case constants.GroupChatType:
				data.ConversationId = data.RecvId
			default:
				_ = srv.Send(websocket.NewErrMessage(errors.New("no selected chat method")), conn)
				return
			}
		}

		err := svc.MsgChatTransferClient.Push(&mq.MsgChatTransfer{
			ConversationId: data.ConversationId,
			ChatType:       data.ChatType,
			SendId:         conn.Uid,
			RecvId:         data.RecvId,
			SendTime:       time.Now().UnixNano(),
			MType:          data.Msg.MType,
			Content:        data.Msg.Content,
		})
		if err != nil {
			_ = srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
	}
}
