package mq

import (
    "github.com/jassue/gin-wire/config"
    "github.com/jassue/gin-wire/utils/path"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
    "path/filepath"
    "time"
)

type QueueLogger struct {
    *zap.Logger
}

func NewQueueLogger(conf *config.Configuration, logger *zap.Logger) *QueueLogger {
    queueWriter := &lumberjack.Logger{
        Filename:   filepath.Join(path.RootPath(), conf.Log.RootDir, "queue.log"),
        MaxSize:    conf.Log.MaxSize,
        MaxBackups: conf.Log.MaxBackups,
        MaxAge:     conf.Log.MaxAge,
        Compress:   conf.Log.Compress,
    }

    var encoder zapcore.Encoder
    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
        encoder.AppendString(time.Format("2006-01-02 15:04:05.000"))
    }
    encoderConfig.EncodeLevel = func(l zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
        encoder.AppendString(conf.App.Env + "." + l.String())
    }
    encoder = zapcore.NewJSONEncoder(encoderConfig)


    o := zap.WrapCore(func(core zapcore.Core) zapcore.Core {
       return zapcore.NewTee(
           zapcore.NewCore(encoder, zapcore.AddSync(queueWriter), zapcore.InfoLevel),
           )
    })

    return &QueueLogger{
        logger.WithOptions(o),
    }
}
