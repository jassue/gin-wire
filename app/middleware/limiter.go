package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/jassue/gin-wire/app/compo"
    cErr "github.com/jassue/gin-wire/app/pkg/error"
    "github.com/jassue/gin-wire/app/pkg/response"
    "golang.org/x/time/rate"
    "time"
)

type Limiter struct {
    lm *compo.LimiterManager
}

func NewLimiterM(lm *compo.LimiterManager) *Limiter {
    return &Limiter{
        lm: lm,
    }
}

func (m *Limiter) Handler(key ...string) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        var limiterKey string
        if len(key) > 0 && len(key[0]) > 0 {
            limiterKey = key[0]
        } else {
            limiterKey = ctx.GetString("token")
            if len(limiterKey) == 0 {
                limiterKey = ctx.ClientIP()
            }
        }

        l := m.lm.GetLimiter(rate.Every(50*time.Millisecond), 300, limiterKey)

        if !l.L.Allow() {
            response.FailByErr(ctx, cErr.TooManyRequestsErr("您的访问过于频繁，请稍候重试"))
            return
        }
    }
}
