package immodels

import (
	"liteChat/pkg/constants"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conversation struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	// TODO: Fill your own fields
	ConversationId string             `bson:"conversation_id,omitempty"`
	ChatType       constants.ChatType `bson:"chat_type,omitempty"`
	IsShow         bool               `bson:"is_show,omitempty"`
	Total          int                `bson:"total,omitempty"`
	Seq            int64              `bson:"seq"`
	Msg            *ChatLog           `bson:"msg,omitempty"`
	UpdateAt       time.Time          `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt       time.Time          `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
