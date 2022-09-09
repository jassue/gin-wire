package handler

import (
    "context"
    "errors"
    "fmt"
    "github.com/jassue/gin-wire/app/compo/mq/event"
)

type ExampleConsumer struct {
}

func NewExampleConsumer() *ExampleConsumer {
    return &ExampleConsumer{
    }
}

func (c *ExampleConsumer) ExampleJob(ctx context.Context, msg event.Event) error {
    e, ok := msg.(*event.ExampleEvent)
    if !ok {
        return errors.New("event type error")
    }
    fmt.Println(e)
    return nil
}
