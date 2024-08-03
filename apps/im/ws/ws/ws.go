package ws

import "liteChat/pkg/constants"

type (
	Msg struct {
		constants.MType `mapstructure:"m_type"`
		Content         string `mapstructure:"content"`
	}

	Chat struct {
		ConversationId     string `mapstructure:"conversation_id"`
		constants.ChatType `mapstructure:"chat_type"`
		SendId             string `mapstructure:"send_id"`
		RecvId             string `mapstructure:"recv_id"`
		SendTime           int64  `mapstructure:"send_time"`
		Msg                `mapstructure:"msg"`
	}

	Push struct {
		ConversationId        string `mapstructure:"conversation_id"`
		constants.ChatType    `mapstructure:"chat_type"`
		SendId                string            `mapstructure:"send_id"`
		RecvId                string            `mapstructure:"recv_id"`
		RecvIds               []string          `mapstructure:"recv_ids"`
		SendTime              int64             `mapstructure:"send_time"`
		ReadRecords           map[string]string `mapstructure:"read_records"`
		constants.ContentType `mapstructure:"content_type"`
		constants.MType       `mapstructure:"m_type"`
		Content               string `mapstructure:"content"`
	}

	MarkRead struct {
		constants.ChatType `mapstructure:"chat_type"`
		RecvId             string   `mapstructure:"recv_id"`
		ConversationId     string   `mapstructure:"conversation_id"`
		MsgIds             []string `mapstructure:"msg_ids"`
	}
)
