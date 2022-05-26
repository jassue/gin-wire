package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/jassue/gin-wire/app/pkg/response"
    "gopkg.in/natefinch/lumberjack.v2"
)

func CustomRecovery(loggerWriter *lumberjack.Logger) gin.HandlerFunc {
    return gin.RecoveryWithWriter(
        loggerWriter,
        response.ServerError)
}
