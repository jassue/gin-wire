package cron

import (
    "context"
    "github.com/google/wire"
    "github.com/jassue/gin-wire/app/data"
    "github.com/robfig/cron/v3"
    "go.uber.org/zap"
)

// ProviderSet is cron providers.
var ProviderSet = wire.NewSet(NewCron, NewExampleJob)

type Cron struct {
    logger *zap.Logger
    data *data.Data
    server *cron.Cron

    exampleJob *ExampleJob
}

// NewCron .
func NewCron(data *data.Data, logger *zap.Logger, exampleJob *ExampleJob) *Cron {
    server := cron.New(
        cron.WithSeconds(),
        )

    return &Cron{
        logger: logger,
        data: data,
        server: server,

        exampleJob: exampleJob,
    }
}

func (c *Cron) Run() error {
    // cron example
    //if _, err := c.server.AddFunc("*/5 * * * * *", c.exampleJob.Hello); err != nil {
    //   return err
    //}

    c.server.Start()
    return nil
}

func (c *Cron) Stop(ctx context.Context) error {
    c.server.Stop()
    return nil
}
