package msgtransfer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"liteChat/apps/im/ws/ws"
	"liteChat/apps/task/mq/internal/svc"
	"liteChat/apps/task/mq/mq"
	"liteChat/pkg/bitmap"
	"liteChat/pkg/constants"
	"sync"
	"time"
)

var (
	GroupMsgReadRecordDelayTime = time.Second
	GroupMsgReadRecordMaxCount  = 10
)

const (
	GroupMsgReadHandlerAtTransfer = iota
	GroupMsgReadHandlerDelayTransfer
)

type MsgReadTransfer struct {
	*BaseMsgTransfer
	cache.Cache
	sync.Mutex
	GroupMsgs map[string]*GroupMsgRead
	pushCh    chan *ws.Push
}

func NewMsgReadTransfer(svc *svc.ServiceContext) kq.ConsumeHandler {
	m := &MsgReadTransfer{
		BaseMsgTransfer: NewBaseMsgTransfer(svc),
		GroupMsgs:       make(map[string]*GroupMsgRead, 1),
		pushCh:          make(chan *ws.Push, 1),
	}

	if svc.Config.MsgReadHandler.GroupMsgReadHandler != GroupMsgReadHandlerAtTransfer {
		if svc.Config.MsgReadHandler.GroupMsgReadRecordMaxCount > 0 {
			GroupMsgReadRecordMaxCount = svc.Config.MsgReadHandler.GroupMsgReadRecordMaxCount
		}

		if svc.Config.MsgReadHandler.GroupMsgReadRecordDelayTime > 0 {
			GroupMsgReadRecordDelayTime = time.Duration(svc.Config.MsgReadHandler.GroupMsgReadRecordDelayTime)
		}
	}

	go m.transfer()

	return m
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

	push := &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ContentType:    constants.ContentMarkRead,
		ReadRecords:    readRecords,
	}

	switch data.ChatType {
	case constants.SingleChatType:
		m.pushCh <- push
	case constants.GroupChatType:
		if m.svcCtx.MsgReadHandler.GroupMsgReadHandler == GroupMsgReadHandlerAtTransfer {
			m.pushCh <- push
		}
		m.Lock()
		defer m.Unlock()
		push.SendId = ""
		if _, ok := m.GroupMsgs[push.ConversationId]; ok {
			m.Infof("Merge push %v", push.ConversationId)
			m.GroupMsgs[push.ConversationId].MergePush(push)
		} else {
			m.Infof("NewGroupMsgRead push %v", push.ConversationId)
			m.GroupMsgs[push.ConversationId] = NewGroupMsgRead(push, m.pushCh)
		}
	}
	return nil
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

func (m *MsgReadTransfer) transfer() {
	for push := range m.pushCh {
		if push.RecvId != "" || len(push.RecvIds) > 0 {
			if err := m.Transfer(context.Background(), push); err != nil {
				m.Errorf("m transfer err %v push %v", err, push)
			}
		}
		if push.ChatType != constants.GroupChatType {
			continue
		}
		if m.svcCtx.MsgReadHandler.GroupMsgReadHandler == GroupMsgReadHandlerAtTransfer {
			continue
		}
		m.Lock()

		if _, ok := m.GroupMsgs[push.ConversationId]; ok && m.GroupMsgs[push.ConversationId].IsActive() {
			m.GroupMsgs[push.ConversationId].Clear()
			delete(m.GroupMsgs, push.ConversationId)
		}
		m.Unlock()
	}
}
