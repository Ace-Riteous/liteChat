package immodels

import (
	"liteChat/pkg/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var DefaultChatLogLimit int64 = 100

type ChatLog struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ConversationId string             `bson:"conversation_id"`
	SendId         string             `bson:"send_id"`
	RecvId         string             `bson:"recv_id"`
	MsgFrom        int                `bson:"msg_from"`
	ChatType       constants.ChatType `bson:"chat_type"`
	MsgType        constants.MType    `bson:"msg_type"`
	MsgContent     string             `bson:"msg_content"`
	SendTime       int64              `bson:"send_time"`
	Status         int                `bson:"status"`
	ReadRecords    []byte             `bson:"read_records"`

	// TODO: Fill your own fields
	UpdateAt time.Time `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt time.Time `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
