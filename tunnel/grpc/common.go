package grpc

import "time"

const (
	DialTimeout     = 5 * time.Second
	ConnectTimeout  = 5 * time.Second
	BackoffMaxDelay = 3 * time.Second

	InitialWindowSize     = 1 << 30
	InitialConnWindowSize = 1 << 30

	KeepAliveTimeout = time.Duration(3) * time.Second
	KeepAliveTime    = time.Duration(10) * time.Second

	MaxRecvMsgSize = 1 << 30
	MaxSendMsgSize = 1 << 30

	ReadBufferSize  = 1 << 10
	WriteBufferSize = 1 << 10

	MaxIdle              = 8
	MaxActive            = 64
	MaxConcurrentStreams = 64
)
