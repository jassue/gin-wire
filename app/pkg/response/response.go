package response

import (
    "github.com/gin-gonic/gin"
    cErr "github.com/jassue/gin-wire/app/pkg/error"
    "net/http"
    "os"
)

type Response struct {
    ErrorCode int `json:"error_code"`
    Data interface{} `json:"data"`
    Message string `json:"message"`
}

func ServerError(c *gin.Context, err interface{}) {
    msg := "Internal Server Error"
    if os.Getenv(gin.EnvGinMode) != gin.ReleaseMode {
        if _, ok := err.(error); ok {
            msg = err.(error).Error()
        }
    }
    FailByErr(c, cErr.InternalServer(msg))
    c.Abort()
}

func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{
        0,
        data,
        "ok",
    })
}

func Fail(c *gin.Context, httpCode int, errorCode int, msg string) {
    c.JSON(httpCode, Response{
        errorCode,
        nil,
        msg,
    })
}

func FailByErr(c *gin.Context, err error) {
    v, ok := err.(*cErr.Error)
    if ok {
        Fail(c, v.HttpCode(), v.ErrorCode(), v.Error())
    } else {
        Fail(c, http.StatusBadRequest, cErr.DEFAULT_ERROR, err.Error())
    }
}
