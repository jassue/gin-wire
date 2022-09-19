package compo

import (
    "golang.org/x/time/rate"
    "sync"
    "time"
)

type LimiterManager struct {
    limiters *sync.Map
    once *sync.Once
}

func NewLimiterManager() *LimiterManager {
    return &LimiterManager{
        limiters: &sync.Map{},
        once: &sync.Once{},
    }
}

type Limiter struct {
    L *rate.Limiter
    lastGetTime time.Time
}

func (lm *LimiterManager) GetLimiter(r rate.Limit, b int, key string) *Limiter {
    lm.once.Do(func() {
        go lm.clearLimiter()
    })

    limiter, ok := lm.limiters.Load(key)
    if ok {
        return limiter.(*Limiter)
    }

    l := &Limiter{
        L: rate.NewLimiter(r, b),
        lastGetTime: time.Now(),
    }

    lm.limiters.Store(key, l)

    return l
}

func (lm *LimiterManager) clearLimiter() {
    for {
        time.Sleep(1 * time.Minute)

        lm.limiters.Range(func(key, value interface{}) bool {
            if time.Now().Unix()-value.(*Limiter).lastGetTime.Unix() > 180 {
                lm.limiters.Delete(key)
            }
            return true
        })
    }
}
