package websocket

import "net/http"

type DialOptions func(opt *DialOption)

type DialOption struct {
	header http.Header
	patten string
}

func NewDialOption(opts ...DialOptions) *DialOption {
	o := DialOption{
		header: nil,
		patten: "/ws",
	}

	for _, opt := range opts {
		opt(&o)
	}

	return &o
}

func WithClientPatten(patten string) DialOptions {
	return func(opt *DialOption) {
		opt.patten = patten
	}
}

func WithClientHeader(header http.Header) DialOptions {
	return func(opt *DialOption) {
		opt.header = header
	}
}
