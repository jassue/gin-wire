package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/jassue/gin-wire/app/pkg/response"
    "gopkg.in/natefinch/lumberjack.v2"
)

type Recovery struct {
    loggerWriter *lumberjack.Logger
}

func NewRecoveryM(loggerWriter *lumberjack.Logger) *Recovery {
    return &Recovery{
        loggerWriter: loggerWriter,
    }
}

func (m *Recovery) Handler() gin.HandlerFunc {
    return gin.RecoveryWithWriter(
        m.loggerWriter,
        response.ServerError,
        )
}
