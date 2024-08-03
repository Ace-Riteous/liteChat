package websocket

import (
	"time"
)

type ServerOptions func(opt *ServerOption)

type ServerOption struct {
	Authentication
	ack               AckType
	ackTimeout        time.Duration
	patten            string
	MaxConnectionIdle time.Duration
	Concurrency       int
}

func NewServerOptions(opts ...ServerOptions) ServerOption {
	o := ServerOption{
		Authentication:    new(authentication),
		patten:            DefaultPatten,
		ack:               NoAck,
		ackTimeout:        DefaultAckTimeout,
		MaxConnectionIdle: DefaultMaxConnectionIdle,
		Concurrency:       DefaultConcurrency,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return o
}

func WithServerAuthentication(auth Authentication) ServerOptions {
	return func(opt *ServerOption) {
		opt.Authentication = auth
	}
}

func WithServerPatten(patten string) ServerOptions {
	return func(opt *ServerOption) {
		opt.patten = patten
	}
}

func WithServerMaxConnectionIdle(maxConnectionIdle time.Duration) ServerOptions {
	return func(opt *ServerOption) {
		if maxConnectionIdle > 0 {
			opt.MaxConnectionIdle = maxConnectionIdle
		}
	}
}

func WithServerAck(ack AckType) ServerOptions {
	return func(opt *ServerOption) {
		opt.ack = ack
	}
}
