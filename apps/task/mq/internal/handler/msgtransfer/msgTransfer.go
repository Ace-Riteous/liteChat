package msgtransfer

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"liteChat/apps/im/ws/websocket"
	"liteChat/apps/im/ws/ws"
	"liteChat/apps/social/rpc/socialclient"
	"liteChat/apps/task/mq/internal/svc"
	"liteChat/pkg/constants"
)

type BaseMsgTransfer struct {
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBaseMsgTransfer(svc *svc.ServiceContext) *BaseMsgTransfer {
	return &BaseMsgTransfer{
		svcCtx: svc,
		Logger: logx.WithContext(context.Background()),
	}
}

func (b *BaseMsgTransfer) Transfer(ctx context.Context, data *ws.Push) error {
	var err error
	switch data.ChatType {
	case constants.GroupChatType:
		err = b.Group(ctx, data)
	case constants.SingleChatType:
		err = b.Single(ctx, data)
	}
	return err
}

func (b *BaseMsgTransfer) Single(ctx context.Context, data *ws.Push) error {
	return b.svcCtx.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FromId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}

func (b *BaseMsgTransfer) Group(ctx context.Context, data *ws.Push) error {
	users, err := b.svcCtx.GroupUsers(ctx, &socialclient.GroupUsersReq{
		GroupId: data.RecvId,
	})
	if err != nil {
		return err
	}
	data.RecvIds = make([]string, len(users.List))
	for _, members := range users.List {
		if members.UserId == data.SendId {
			continue
		}
		data.RecvIds = append(data.RecvIds, members.UserId)
	}

	return b.Single(ctx, data)
}
