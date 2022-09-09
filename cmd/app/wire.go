//go:build wireinject
// +build wireinject

package main

import (
    "github.com/google/wire"
    "github.com/jassue/gin-wire/app/command"
    commandH "github.com/jassue/gin-wire/app/command/handler"
    "github.com/jassue/gin-wire/app/compo"
    "github.com/jassue/gin-wire/app/cron"
    "github.com/jassue/gin-wire/app/data"
    "github.com/jassue/gin-wire/app/handler"
    "github.com/jassue/gin-wire/app/middleware"
    "github.com/jassue/gin-wire/app/service"
    "github.com/jassue/gin-wire/config"
    "github.com/jassue/gin-wire/router"
    "go.uber.org/zap"
    "gopkg.in/natefinch/lumberjack.v2"
)

// wireApp init application.
func wireApp(*config.Configuration, *lumberjack.Logger, *zap.Logger) (*App, func(), error) {
    panic(
        wire.Build(
            data.ProviderSet,
            compo.ProviderSet,
            service.ProviderSet,
            handler.ProviderSet,
            middleware.ProviderSet,
            router.ProviderSet,
            cron.ProviderSet,
            //consumer.ProviderSet,
            newHttpServer,
            newApp,
            ),
        )
}

// wireCommand init application.
func wireCommand(*config.Configuration, *lumberjack.Logger, *zap.Logger) (*command.Command, func(), error) {
   panic(wire.Build(commandH.ProviderSet, command.NewCommand))
}

