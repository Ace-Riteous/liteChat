package user

import (
	"liteChat/apps/im/ws/internal/svc"
	"liteChat/apps/im/ws/websocket"
)

func Online(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		uids := srv.GetUids()
		err := srv.Send(websocket.NewMessage(srv.GetUid(conn), uids), conn)
		srv.Info("err ", err)
	}
}
