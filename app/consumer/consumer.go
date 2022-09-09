package consumer

import (
    "context"
    "github.com/google/wire"
    "github.com/jassue/gin-wire/app/compo/mq/event"
    "github.com/jassue/gin-wire/app/consumer/handler"
)

// ProviderSet is cron providers.
var ProviderSet = wire.NewSet(NewConsumer, handler.NewExampleConsumer)

type Consumer struct {
    mqReceiver event.Receiver
    cancelFunc context.CancelFunc
    exampleC *handler.ExampleConsumer
}

func NewConsumer(r event.Receiver, exampleC *handler.ExampleConsumer) *Consumer {
    return &Consumer{
        mqReceiver: r,
        exampleC: exampleC,
    }
}

func (c *Consumer) Run() error {
    ctx, cancelFunc := context.WithCancel(context.Background())
    c.cancelFunc = cancelFunc

    _ = c.mqReceiver.Receive(ctx, &event.ExampleEvent{}, c.exampleC.ExampleJob, 3)

    return nil
}

func (c *Consumer) Stop(ctx context.Context) error {
    c.cancelFunc()
    _ = c.mqReceiver.Close()
    return nil
}
