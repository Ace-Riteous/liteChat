package immodels

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conversations struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	// TODO: Fill your own fields
	UserId           string                   `bson:"user_id"`
	ConversationList map[string]*Conversation `bson:"conversation_list"`
	UpdateAt         time.Time                `bson:"updateAt,omitempty" json:"updateAt,omitempty"`
	CreateAt         time.Time                `bson:"createAt,omitempty" json:"createAt,omitempty"`
}
