package msgtransfer

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"liteChat/apps/im/immodels"
	"liteChat/apps/im/ws/ws"
	"liteChat/apps/task/mq/internal/svc"
	"liteChat/apps/task/mq/mq"
	"liteChat/pkg/bitmap"
)

type MsgChatTransfer struct {
	*BaseMsgTransfer
}

func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		NewBaseMsgTransfer(svc),
	}
}

func (m *MsgChatTransfer) Consume(ctx context.Context, key, value string) error {
	fmt.Println("key : ", key, "value : ", value)

	var (
		data  mq.MsgChatTransfer
		msgId = primitive.NewObjectID()
	)
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	if err := m.AddChatLog(ctx, msgId, &data); err != nil {
		return err
	}

	return m.Transfer(ctx, &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		RecvIds:        data.RecvIds,
		SendTime:       data.SendTime,
		MsgId:          msgId.Hex(),
		MType:          data.MType,
		Content:        data.Content,
	})
}

func (m *MsgChatTransfer) AddChatLog(ctx context.Context, msgId primitive.ObjectID, data *mq.MsgChatTransfer) error {
	chatLog := immodels.ChatLog{
		ID:             msgId,
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ChatType:       data.ChatType,
		MsgFrom:        0,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       data.SendTime,
	}

	//将自己设置为已读
	readRecords := bitmap.NewBitmap(0)
	readRecords.Set(chatLog.SendId)
	chatLog.ReadRecords = readRecords.Export()

	err := m.svcCtx.ChatLogModel.Insert(ctx, &chatLog)
	if err != nil {
		return err
	}
	return m.svcCtx.ConversationModel.UpdateMsg(ctx, &chatLog)
}
