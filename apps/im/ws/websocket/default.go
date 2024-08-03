package websocket

import (
	"time"
)

const (
	DefaultPatten            = "/ws"
	DefaultMaxConnectionIdle = 2 * 60 * time.Second
	DefaultAckTimeout        = 30 * time.Second
	DefaultConcurrency       = 10
)
