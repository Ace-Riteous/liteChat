package msgtransfer

import (
	"github.com/zeromicro/go-zero/core/logx"
	"liteChat/apps/im/ws/ws"
	"liteChat/pkg/constants"
	"sync"
	"time"
)

type GroupMsgRead struct {
	sync.Mutex
	ConversationId string
	Push           *ws.Push
	PushCh         chan *ws.Push
	Count          int
	PushTime       time.Time
	Done           chan struct{}
}

func NewGroupMsgRead(push *ws.Push, pushCh chan *ws.Push) *GroupMsgRead {
	g := &GroupMsgRead{
		ConversationId: push.ConversationId,
		Push:           push,
		PushCh:         pushCh,
		Count:          1,
		PushTime:       time.Now(),
		Done:           make(chan struct{}),
	}

	go g.transfer()

	return g
}

func (g *GroupMsgRead) MergePush(push *ws.Push) {
	g.Lock()
	defer g.Unlock()
	g.Count++
	for msgId, read := range push.ReadRecords {
		g.Push.ReadRecords[msgId] = read
	}
}

func (g *GroupMsgRead) transfer() {
	//超时发送
	//超量发送
	timer := time.NewTimer(GroupMsgReadRecordDelayTime / 2)
	defer timer.Stop()
	for {
		select {
		case <-g.Done:
			return
		case <-timer.C:
			g.Lock()
			pushTime := g.PushTime
			val := GroupMsgReadRecordDelayTime*2 - time.Since(pushTime)
			push := g.Push

			if val > 0 && g.Count < GroupMsgReadRecordMaxCount || push == nil {
				if val > 0 {
					timer.Reset(val)
				}

				g.Unlock()
				continue
			}

			g.PushTime = time.Now()
			g.Push = nil
			g.Count = 0
			timer.Reset(GroupMsgReadRecordDelayTime / 2)
			g.Unlock()

			logx.Infof("超过合并条件，推送 %v", push)
			g.PushCh <- push
		default:
			g.Lock()
			if g.Count >= GroupMsgReadRecordMaxCount {
				push := g.Push
				g.Push = nil
				g.Count = 0
				g.Unlock()

				logx.Infof("default 超过合并条件，推送 %v", push)
				g.PushCh <- push
				continue
			}

			if g.isActive() {
				g.Unlock()
				//释放msgReadTransfer
				push := &ws.Push{
					ConversationId: g.ConversationId,
					ChatType:       constants.GroupChatType,
				}
				g.PushCh <- push
				continue
			}
			g.Unlock()
			tempDelay := GroupMsgReadRecordDelayTime / 4
			if tempDelay > time.Second {
				tempDelay = time.Second
			}
			time.Sleep(tempDelay)
		}
	}

}

func (g *GroupMsgRead) IsActive() bool {
	g.Lock()
	defer g.Unlock()
	return g.isActive()
}

func (g *GroupMsgRead) isActive() bool {
	pushTime := g.PushTime
	val := GroupMsgReadRecordDelayTime*2 - time.Since(pushTime)
	if val <= 0 && g.PushCh == nil && g.Count == 0 {
		return true
	}
	return false
}

func (g *GroupMsgRead) Clear() {
	select {
	case <-g.Done:
		return
	default:
		close(g.Done)
	}
	g.PushCh = nil
}
