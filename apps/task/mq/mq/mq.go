package mq

import "liteChat/pkg/constants"

type MsgChatTransfer struct {
	ConversationId     string `json:"conversation_id"`
	constants.ChatType `json:"chat_type"`
	SendId             string   `json:"send_id"`
	RecvId             string   `json:"recv_id"`
	RecvIds            []string `json:"recv_ids"`
	SendTime           int64    `json:"send_time"`

	constants.MType `json:"m_type"`
	Content         string `json:"content"`
}

type MsgMarkRead struct {
	ConversationId     string `json:"conversation_id"`
	constants.ChatType `json:"chat_type"`
	SendId             string   `json:"send_id"`
	RecvId             string   `json:"recv_id"`
	MsgIds             []string `json:"msg_ids"`
}
