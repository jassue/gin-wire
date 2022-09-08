package event

import (
    "context"
)

// Event 消息事件
type Event interface {
    Queue() string
    Value() []byte
    DelaySeconds() int64
    Unmarshal(value []byte) Event
}

// Handler 消息处理
type Handler func(context.Context, Event) error

// Sender 消息生产者
type Sender interface {
    Send(context.Context, Event) error
    Close() error
}

// Receiver 消息接收者
type Receiver interface {
    Receive(ctx context.Context, msg Event, handler Handler, workerNum int64) error
    Close() error
}
