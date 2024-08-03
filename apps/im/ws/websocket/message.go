package websocket

import "time"

type FrameType uint8

const (
	FrameData FrameType = 0x0
	//FrameHeaders      FrameType = 0x1
	//FramePriority     FrameType = 0x2
	//FrameRSTStream    FrameType = 0x3
	//FrameSettings     FrameType = 0x4
	//FramePushPromise  FrameType = 0x5
	FramePing FrameType = 0x1
	//FrameGoAway       FrameType = 0x7
	//FrameWindowUpdate FrameType = 0x8
	//FrameContinuation FrameType = 0x9
	FrameAck   FrameType = 0x2
	FrameNoAck FrameType = 0x3
	FrameErr   FrameType = 0x9
)

type Message struct {
	FrameType `json:"frame_type"`
	Id        string      `json:"id"`
	AckSeq    int         `json:"ack_seq"`
	AckTime   time.Time   `json:"ack_time"`
	ErrCount  int         `json:"err_count"`
	Method    string      `json:"method"`
	FromId    string      `json:"from_id"`
	Data      interface{} `json:"data"`
}

func NewMessage(fromId string, data interface{}) *Message {
	return &Message{
		FrameType: FrameData,
		FromId:    fromId,
		Data:      data,
	}
}

func NewErrMessage(err error) *Message {
	return &Message{
		FrameType: FrameErr,
		Data:      err.Error(),
	}
}
