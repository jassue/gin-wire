package main

import (
    "context"
    "fmt"
    "github.com/fsnotify/fsnotify"
    "github.com/gin-gonic/gin/binding"
    "github.com/go-playground/validator/v10"
    "github.com/jassue/gin-wire/app/command"
    "github.com/jassue/gin-wire/config"
    "github.com/jassue/gin-wire/utils"
    "github.com/spf13/cobra"
    "github.com/spf13/pflag"
    "github.com/spf13/viper"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
    "log"
    "os"
    "os/signal"
    "path/filepath"
    "reflect"
    "strings"
    "syscall"
    "time"
)

var (
    rootPath = utils.RootPath()

    Version string
    configPath string
    conf *config.Configuration
    loggerWriter *lumberjack.Logger
    logger *zap.Logger
)

func init() {
    pflag.StringVarP(&configPath, "conf", "", filepath.Join(rootPath, "conf", "config.yaml"), "config path, eg: --conf config.yaml")

    cobra.OnInitialize(func() {
        initConfig()
        initLogger()
        initValidator()
    })
}

func main() {
    rootCmd := &cobra.Command{
        Use:   "app",
        Run: func(cmd *cobra.Command, args []string) {
            app, cleanup, err := wireApp(conf, loggerWriter, logger)
            if err != nil {
                panic(err)
            }
            defer cleanup()

            // 启动应用
            log.Printf("start app %s ...", Version)
            if err := app.Run(); err != nil {
                panic(err)
            }

            // 等待中断信号以优雅地关闭应用
            quit := make(chan os.Signal)
            signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
            <-quit

            log.Printf("shutdown app %s ...", Version)

            // 设置 5 秒的超时时间
            ctx, cancel := context.WithTimeout(app.cxt, 5*time.Second)
            defer cancel()

            // 关闭应用
            if err := app.Stop(ctx); err != nil {
                panic(err)
            }
        },
    }

    // 注册命令
    command.Register(rootCmd, func() (*command.Command, func(), error) {
       return wireCommand(conf, loggerWriter, logger)
    })

    if err := rootCmd.Execute(); err != nil {
        panic(err)
    }
}

func initConfig() {
    if !filepath.IsAbs(configPath) {
        configPath = filepath.Join(rootPath, "conf", configPath)
    }

    fmt.Println("load config:" + configPath)

    v := viper.New()
    v.SetConfigFile(configPath)
    v.SetConfigType("yaml")
    if err := v.ReadInConfig(); err != nil {
        panic(fmt.Errorf("read config failed: %s \n", err))
    }

    if err := v.Unmarshal(&conf); err != nil {
        fmt.Println(err)
    }

    v.WatchConfig()
    v.OnConfigChange(func(in fsnotify.Event) {
        fmt.Println("config file changed:", in.Name)
        defer func() {
            if err := recover(); err != nil {
                logger.Error("config file changed err:", zap.Any("err", err))
                fmt.Println(err)
            }
        }()
        if err := v.Unmarshal(&conf); err != nil {
            fmt.Println(err)
        }
    })
}

func initLogger() {
    var level zapcore.Level // zap 日志等级
    var options []zap.Option // zap 配置项

    logFileDir := conf.Log.RootDir
    if !filepath.IsAbs(logFileDir) {
        logFileDir = filepath.Join(rootPath, logFileDir)
    }

    if ok, _ := utils.PathExists(logFileDir); !ok {
        _ = os.Mkdir(conf.Log.RootDir, os.ModePerm)
    }

    switch conf.Log.Level {
    case "debug":
        level = zap.DebugLevel
        options = append(options, zap.AddStacktrace(level))
    case "info":
        level = zap.InfoLevel
    case "warn":
        level = zap.WarnLevel
    case "error":
        level = zap.ErrorLevel
        options = append(options, zap.AddStacktrace(level))
    case "dpanic":
        level = zap.DPanicLevel
    case "panic":
        level = zap.PanicLevel
    case "fatal":
        level = zap.FatalLevel
    default:
        level = zap.InfoLevel
    }

    // 调整编码器默认配置
    var encoder zapcore.Encoder
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
        encoder.AppendString(time.Format("[" + "2006-01-02 15:04:05.000" + "]"))
    }
    encoderConfig.EncodeLevel = func(l zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
        encoder.AppendString(conf.App.Env + "." + l.String())
    }

    // 设置编码器
    if conf.Log.Format == "json" {
        encoder = zapcore.NewJSONEncoder(encoderConfig)
    } else {
        encoder = zapcore.NewConsoleEncoder(encoderConfig)
    }

    loggerWriter = &lumberjack.Logger{
        Filename:   filepath.Join(logFileDir, conf.Log.Filename),
        MaxSize:    conf.Log.MaxSize,
        MaxBackups: conf.Log.MaxBackups,
        MaxAge:     conf.Log.MaxAge,
        Compress:   conf.Log.Compress,
    }

    logger = zap.New(zapcore.NewCore(encoder, zapcore.AddSync(loggerWriter), level), options...)
}

func initValidator() {
    if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
        // 注册自定义验证器
        _ = v.RegisterValidation("mobile", utils.ValidateMobile)

        // 注册自定义 json tag 函数
        v.RegisterTagNameFunc(func(fld reflect.StructField) string {
            name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
            if name == "-" {
                return ""
            }
            return name
        })
    }
}
