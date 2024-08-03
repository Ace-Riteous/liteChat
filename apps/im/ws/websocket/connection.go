package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type Conn struct {
	*websocket.Conn
	IdleMu            sync.Mutex
	Uid               string
	S                 *Server
	MsgMu             sync.Mutex
	readMessage       []*Message
	readMessageSeq    map[string]*Message
	MessageChan       chan *Message
	Idle              time.Time
	MaxConnectionIdle time.Duration
	Done              chan struct{}
}

func NewConn(s *Server, w http.ResponseWriter, r *http.Request) *Conn {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("upgrade err %v", err)
		return nil
	}

	conn := &Conn{
		Conn:              c,
		S:                 s,
		Idle:              time.Now(),
		MaxConnectionIdle: s.opt.MaxConnectionIdle,
		readMessage:       make([]*Message, 0, 2),
		readMessageSeq:    make(map[string]*Message, 2),
		MessageChan:       make(chan *Message, 1),
		Done:              make(chan struct{}),
	}

	go conn.keepalive()

	return conn
}

func (c *Conn) AppendMsgMq(msg *Message) {
	c.MsgMu.Lock()
	defer c.MsgMu.Unlock()

	if m, ok := c.readMessageSeq[msg.Id]; ok {
		if len(c.readMessage) == 0 {
			return
		}
		if msg.AckSeq <= m.AckSeq {
			//还没有ack确认,或者重复
			return
		}
		c.readMessageSeq[msg.Id] = msg
		return
	}
	//避免客户端重复ack验证
	if msg.FrameType == FrameAck {
		return
	}
	c.readMessage = append(c.readMessage, msg)
	c.readMessageSeq[msg.Id] = msg
}
func (c *Conn) ReadMessage() (messageType int, p []byte, err error) {
	c.IdleMu.Lock()
	defer c.IdleMu.Unlock()
	messageType, p, err = c.Conn.ReadMessage()
	c.Idle = time.Time{}
	return messageType, p, err
}

func (c *Conn) WriteMessage(messageType int, data []byte) error {
	c.IdleMu.Lock()
	defer c.IdleMu.Unlock()
	err := c.Conn.WriteMessage(messageType, data)
	c.Idle = time.Now()
	return err
}

func (c *Conn) Close() error {
	select {
	case <-c.Done:

	default:
		close(c.Done)
	}

	return c.Conn.Close()
}

func (c *Conn) keepalive() {
	idleTimer := time.NewTimer(c.MaxConnectionIdle)
	defer func() {
		idleTimer.Stop()
	}()

	for {
		select {
		case <-idleTimer.C:
			c.IdleMu.Lock()

			if c.Idle.IsZero() {
				c.IdleMu.Unlock()
				idleTimer.Reset(c.MaxConnectionIdle)
				continue
			}
			val := c.MaxConnectionIdle - time.Since(c.Idle)
			c.IdleMu.Unlock()
			if val <= 0 {
				c.S.Close(c)
				return
			}
			idleTimer.Reset(val)
		case <-c.Done:
			return
		}
	}
}
