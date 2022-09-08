package compo

import (
    "github.com/google/wire"
    "github.com/jassue/gin-wire/app/compo/casbin"
    "github.com/jassue/gin-wire/app/compo/mq"
    "github.com/jassue/gin-wire/app/compo/mq/rabbitmq"
)

// ProviderSet is compo providers.
var ProviderSet = wire.NewSet(
    NewSonyFlake,
    NewLockBuilder,
    NewStorage,
    casbin.NewEnforcer,
    mq.NewQueueLogger,
    rabbitmq.NewConnManager,
    rabbitmq.NewRabbitmqSender,
    rabbitmq.NewRabbitmqReceiver,
    )
