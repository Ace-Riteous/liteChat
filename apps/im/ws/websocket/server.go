package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"net/http"
	"sync"
	"time"
)

type AckType int

const (
	NoAck = iota
	OnlyAck
	RigorAck
)

func (t AckType) ToString() string {
	switch t {
	case OnlyAck:
		return "OnlyAck"
	case RigorAck:
		return "RigorAck"
	default:
		return "NoAck"
	}
}

type Server struct {
	sync.RWMutex
	opt *ServerOption
	*threading.TaskRunner
	authentication Authentication
	patten         string
	routes         map[string]HandlerFunc
	addr           string
	connToUser     map[*Conn]string
	userToConn     map[string]*Conn
	upgrader       websocket.Upgrader
	logx.Logger
}

func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := NewServerOptions(opts...)
	return &Server{
		opt:            &opt,
		TaskRunner:     threading.NewTaskRunner(opt.Concurrency),
		authentication: opt.Authentication,
		patten:         opt.patten,
		routes:         make(map[string]HandlerFunc),
		addr:           addr,
		upgrader:       websocket.Upgrader{},
		connToUser:     make(map[*Conn]string),
		userToConn:     make(map[string]*Conn),
		Logger:         logx.WithContext(context.Background()),
	}
}

func (s *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			s.Errorf("service recovered from err %v", r)
		}
	}()

	conn := NewConn(s, w, r)
	if conn == nil {
		return
	}

	if !s.authentication.Auth(w, r) {
		//_ = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint("no enough auth")))
		_ = s.Send(&Message{
			FrameType: FrameData,
			Data:      fmt.Sprint("no enough auth"),
		}, conn)
		_ = conn.Close()
		return
	}
	s.AddConn(conn, r)

	go s.HandlerConn(conn)
}

func (s *Server) HandlerConn(conn *Conn) {
	uids := s.GetUids(conn)
	conn.Uid = uids[0]

	//处理任务
	go s.HandlerWrite(conn)

	if s.IsAck(nil) {
		go s.ReadAck(conn)
	}
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("ws read msg err %v", err)
			s.Close(conn)
			return
		}
		var message Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			s.Errorf("json unmarshal err %v", err)
			s.Close(conn)
			return
		}

		if s.IsAck(&message) {
			s.Infof("conn read ack msg %v", message)
			conn.AppendMsgMq(&message)
		} else {
			conn.MessageChan <- &message
		}
	}
}

func (s *Server) IsAck(message *Message) bool {
	if message == nil {
		return s.opt.ack != NoAck
	}
	return s.opt.ack != NoAck && message.FrameType != FrameNoAck

}
func (s *Server) ReadAck(conn *Conn) {
	send := func(message *Message, conn *Conn) error {
		err := s.Send(message, conn)
		if err == nil {
			return nil
		}
		s.Errorf("msg ack OnlyAck send err %v msg %v", err, message)
		conn.readMessage[0].ErrCount++
		conn.MsgMu.Unlock()

		tempDelay := time.Duration(200*conn.readMessage[0].ErrCount) * time.Microsecond
		if maxDelay := 1 * time.Second; tempDelay > maxDelay {
			tempDelay = maxDelay
		}
		time.Sleep(tempDelay)
		return err
	}
	for {
		select {
		case <-conn.Done:
			s.Infof("close msg ack uid %v", conn.Uid)
			return
		default:
		}

		conn.MsgMu.Lock()
		if len(conn.readMessage) == 0 {
			conn.MsgMu.Unlock()
			//增加睡眠，有助于任务切换
			time.Sleep(100 * time.Millisecond)
			continue
		}
		//读取第一条
		message := conn.readMessage[0]
		switch s.opt.ack {
		case OnlyAck:
			_ = send(&Message{
				FrameType: FrameAck,
				Id:        message.Id,
				AckSeq:    message.AckSeq + 1,
			}, conn)
			//移除消息
			conn.readMessage = conn.readMessage[1:]
			conn.MsgMu.Unlock()
			conn.MessageChan <- message
		case RigorAck:
			if message.AckSeq == 0 {
				conn.readMessage[0].AckSeq++
				conn.readMessage[0].AckTime = time.Now()
				_ = send(&Message{
					FrameType: FrameAck,
					Id:        message.Id,
					AckSeq:    message.AckSeq,
				}, conn)
				s.Infof("msg ack RigorAck send id %v, seq %v, time %v", message.Id, message.AckSeq, message.AckTime)
				conn.MsgMu.Unlock()
				continue
			}
			//1.客户端返回确认
			msgSeq := conn.readMessageSeq[message.Id]
			if msgSeq.AckSeq > message.AckSeq {
				conn.readMessage = conn.readMessage[1:]
				conn.MsgMu.Unlock()
				conn.MessageChan <- message
				s.Infof("msg ack RigorAck success mid %v", message.Id)
				continue
			}
			//2.客户端未返回
			val := s.opt.ackTimeout - time.Since(message.AckTime)
			if !message.AckTime.IsZero() && val > 0 {
				//  2.1未超时
				conn.MsgMu.Unlock()
				_ = send(&Message{
					FrameType: FrameAck,
					Id:        message.Id,
					AckSeq:    message.AckSeq,
				}, conn)
				time.Sleep(300 * time.Millisecond)
			} else {
				//  2.2超时
				delete(conn.readMessageSeq, message.Id)
				conn.readMessage = conn.readMessage[1:]
				conn.MsgMu.Unlock()
				continue
			}
		default:
			panic("no such ack handler")
		}
	}
}

func (s *Server) HandlerWrite(conn *Conn) {
	for {
		select {
		case <-conn.Done:
			return
		case message := <-conn.MessageChan:
			switch message.FrameType {
			case FramePing:
				_ = s.Send(&Message{
					FrameType: FramePing,
				}, conn)
			case FrameData:
				if handler, ok := s.routes[message.Method]; ok {
					handler(s, conn, message)
				} else {
					_ = s.Send(&Message{
						FrameType: FrameData,
						Data:      fmt.Sprintf("no such hanler, method :%v", message.Method),
					}, conn)
				}
			}

			if s.IsAck(message) {
				conn.MsgMu.Lock()
				delete(conn.readMessageSeq, message.Id)
				conn.MsgMu.Unlock()
			}
		}
	}
}

func (s *Server) AddConn(conn *Conn, r *http.Request) {
	uid := s.authentication.UserId(r)
	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	if c := s.userToConn[uid]; c != nil {
		//不支持重复登陆，踢掉之前的连接
		_ = c.Close()
	}
	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
}

func (s *Server) GetConn(uid string) *Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	return s.userToConn[uid]
}

func (s *Server) GetConns(uids ...string) []*Conn {
	if len(uids) == 0 {
		return nil
	}
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	res := make([]*Conn, 0, len(uids))
	for _, uid := range uids {
		res = append(res, s.userToConn[uid])
	}
	return res
}

func (s *Server) GetUid(conn *Conn) string {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()
	return s.connToUser[conn]
}

func (s *Server) GetUids(conns ...*Conn) []string {

	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	var res []string
	if len(conns) == 0 {
		res = make([]string, 0, len(s.connToUser))
		for _, uid := range s.connToUser {
			res = append(res, uid)
		}
	} else {
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			res = append(res, s.connToUser[conn])
		}
	}

	return res
}

func (s *Server) AddRoutes(rs []Route) {
	for _, r := range rs {
		s.routes[r.Method] = r.Handler
	}
}

func (s *Server) SendByUserId(msg interface{}, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}
	return s.Send(msg, s.GetConns(sendIds...)...)
}

func (s *Server) Send(msg interface{}, conns ...*Conn) error {
	if len(conns) == 0 {
		return nil
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	for _, conn := range conns {
		if err = conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) Start() {
	http.HandleFunc(s.patten, s.ServerWs)
	s.Info(http.ListenAndServe(s.addr, nil))
}

func (s *Server) Stop() {
	fmt.Println("stop server")
}

func (s *Server) Close(conn *Conn) {

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	uid := s.connToUser[conn]
	if uid == "" {
		//已经被关闭
		return
	}
	delete(s.userToConn, uid)
	delete(s.connToUser, conn)
	_ = conn.Close()
}
