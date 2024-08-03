package msgtransfer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
	"liteChat/apps/im/ws/ws"
	"liteChat/apps/task/mq/internal/svc"
	"liteChat/apps/task/mq/mq"
	"liteChat/pkg/bitmap"
	"liteChat/pkg/constants"
)

type MsgReadTransfer struct {
	*BaseMsgTransfer
}

func NewMsgReadTransfer(svc *svc.ServiceContext) kq.ConsumeHandler {
	return &MsgReadTransfer{
		NewBaseMsgTransfer(svc),
	}
}

func (m *MsgReadTransfer) Consume(ctx context.Context, key, value string) error {
	m.Info("MsgReadTransfer ", value)
	var (
		data mq.MsgMarkRead
	)
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	readRecords, err := m.UpdateChatLogRead(ctx, &data)
	if err != nil {
		return err
	}

	return m.Transfer(ctx, &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ContentType:    constants.ContentMarkRead,
		ReadRecords:    readRecords,
	})
}

func (m *MsgReadTransfer) UpdateChatLogRead(ctx context.Context, data *mq.MsgMarkRead) (map[string]string, error) {
	res := make(map[string]string)

	chatLogs, err := m.svcCtx.ChatLogModel.ListByMsgIds(ctx, data.MsgIds)
	if err != nil {
		return nil, err
	}

	for _, chatLog := range chatLogs {
		switch chatLog.ChatType {
		case constants.SingleChatType:
			chatLog.ReadRecords = []byte{1}
		case constants.GroupChatType:
			readRecords := bitmap.Load(chatLog.ReadRecords)
			readRecords.Set(data.SendId)
			chatLog.ReadRecords = readRecords.Export()
		}
		res[chatLog.ID.Hex()] = base64.StdEncoding.EncodeToString(chatLog.ReadRecords)
		err = m.svcCtx.ChatLogModel.UpdateMarkRead(ctx, chatLog.ID, chatLog.ReadRecords)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}
