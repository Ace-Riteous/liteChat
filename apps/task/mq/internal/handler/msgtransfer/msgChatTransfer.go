package msgtransfer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"liteChat/apps/im/immodels"
	"liteChat/apps/im/ws/websocket"
	"liteChat/apps/social/rpc/socialclient"
	"liteChat/apps/task/mq/internal/svc"
	"liteChat/apps/task/mq/mq"
	"liteChat/pkg/constants"
	"liteChat/pkg/xerr"
)

type MsgChatTransfer struct {
	logx.Logger
	svc *svc.ServiceContext
}

func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		Logger: logx.WithContext(context.Background()),
		svc:    svc,
	}
}

func (m *MsgChatTransfer) Consume(ctx context.Context, key, value string) error {
	fmt.Println("key : ", key, "value : ", value)

	var (
		data mq.MsgChatTransfer
	)
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	if err := m.AddChatLog(ctx, &data); err != nil {
		return err
	}

	switch data.ChatType {
	case constants.SingleChatType:
		return m.Single(&data)
	case constants.GroupChatType:
		return m.Group(ctx, &data)
	default:
		return errors.Wrap(xerr.NewInternalErr(), "no such chat_type")
	}
}

func (m *MsgChatTransfer) Single(data *mq.MsgChatTransfer) error {
	return m.svc.WsClient.Send(websocket.Message{
		FrameType: websocket.FrameData,
		Method:    "push",
		FromId:    constants.SYSTEM_ROOT_UID,
		Data:      data,
	})
}

func (m *MsgChatTransfer) Group(ctx context.Context, data *mq.MsgChatTransfer) error {
	users, err := m.svc.GroupUsers(ctx, &socialclient.GroupUsersReq{
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

	return m.Single(data)
}

func (m *MsgChatTransfer) AddChatLog(ctx context.Context, data *mq.MsgChatTransfer) error {
	chatLog := immodels.ChatLog{
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ChatType:       data.ChatType,
		MsgFrom:        0,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       data.SendTime,
	}
	err := m.svc.ChatLogModel.Insert(ctx, &chatLog)
	if err != nil {
		return err
	}
	return m.svc.ConversationModel.UpdateMsg(ctx, &chatLog)
}
