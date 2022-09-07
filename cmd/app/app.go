package main

import (
    "context"
    "github.com/gin-gonic/gin"
    "github.com/jassue/gin-wire/app/cron"
    "github.com/jassue/gin-wire/config"
    "go.uber.org/zap"
    "net/http"
)

type App struct {
    conf *config.Configuration
    logger *zap.Logger
    httpSrv *http.Server
    cronSrv *cron.Cron
}

func newHttpServer(
    conf *config.Configuration,
    router *gin.Engine,
    ) *http.Server {
    return &http.Server{
        Addr:    ":" + conf.App.Port,
        Handler: router,
    }
}

func newApp(
    conf *config.Configuration,
    logger *zap.Logger,
    httpSrv *http.Server,
    cronSrv *cron.Cron,
) *App {
    return &App{
        conf: conf,
        logger: logger,
        httpSrv: httpSrv,
        cronSrv: cronSrv,
    }
}

func (a *App) Run() error {
    // 启动 http server
    go func() {
        a.logger.Info("http server started")
        if err := a.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            panic(err)
        }
    }()

    // 启动 cron server
    go func() {
        a.logger.Info("cron server started")
        if err := a.cronSrv.Run(); err != nil {
            panic(err)
        }
    }()

    return nil
}

func (a *App) Stop(ctx context.Context) error {
    // 关闭 http server
    a.logger.Info("http server has been stop")
    if err := a.httpSrv.Shutdown(ctx); err != nil {
        return err
    }

    // 关闭 cron server
    a.logger.Info("cron server has been stop")
    if err := a.cronSrv.Stop(ctx); err != nil {
        return err
    }

    return nil
}
