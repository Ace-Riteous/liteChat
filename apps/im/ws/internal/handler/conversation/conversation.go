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
			SendTime:       time.Now().UnixMilli(),
			MType:          data.Msg.MType,
			Content:        data.Msg.Content,
			MsgId:          msg.Id,
		})
		if err != nil {
			_ = srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
	}
}

func MarkRead(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var data ws.MarkRead
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			_ = srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
		err := svc.MsgReadTransferClient.Push(&mq.MsgMarkRead{
			ConversationId: data.ConversationId,
			ChatType:       data.ChatType,
			SendId:         conn.Uid,
			RecvId:         data.RecvId,
			MsgIds:         data.MsgIds,
		})
		if err != nil {
			_ = srv.Send(websocket.NewErrMessage(err), conn)
			return
		}
	}
}
