package main

import (
    "context"
    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
    "github.com/go-playground/validator/v10"
    "github.com/jassue/gin-wire/app/cron"
    "github.com/jassue/gin-wire/config"
    validator2 "github.com/jassue/gin-wire/utils/validator"
    "go.uber.org/zap"
    "net/http"
    "reflect"
    "strings"
)

type App struct {
    conf *config.Configuration
    logger *zap.Logger
    httpSrv *http.Server
    cronSrv *cron.Cron
    //consumerSrv *consumer.Consumer
}

func newHttpServer(
    conf *config.Configuration,
    router *gin.Engine,
    ) *http.Server {
    initValidator()
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
    //consumerSrv *consumer.Consumer,
) *App {
    return &App{
        conf: conf,
        logger: logger,
        httpSrv: httpSrv,
        cronSrv: cronSrv,
        //consumerSrv: consumerSrv,
    }
}

func initValidator() {
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        // 注册自定义验证器
        _ = v.RegisterValidation("mobile", validator2.ValidateMobile)

        // 注册自定义 tag 函数
        v.RegisterTagNameFunc(func(fld reflect.StructField) string {
            // 'vn' tag - ValidatorMessages key name
            name := strings.SplitN(fld.Tag.Get("vn"), ",", 2)[0]
            if name == "-" {
                return ""
            }
            return name
        })
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

    // 启动 queue consumer
    //go func() {
    //   a.logger.Info("queue worker started")
    //   if err := a.consumerSrv.Run(); err != nil {
    //       panic(err)
    //   }
    //}()

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

    // 关闭 queue consumer
    //a.logger.Info("queue consumer has been stop")
    //if err := a.consumerSrv.Stop(ctx); err != nil {
    //   return err
    //}

    return nil
}
