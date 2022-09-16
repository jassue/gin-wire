package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt"
    "github.com/jassue/gin-wire/app/domain"
    cErr "github.com/jassue/gin-wire/app/pkg/error"
    "github.com/jassue/gin-wire/app/pkg/response"
    "github.com/jassue/gin-wire/app/service"
    "github.com/jassue/gin-wire/config"
    "strconv"
    "time"
)

type JWTAuth struct {
    conf *config.Configuration
    jwtS *service.JwtService
}

func NewJWTAuthM(
    conf *config.Configuration,
    jwtS *service.JwtService,
    ) *JWTAuth {
    return &JWTAuth{
        conf: conf,
        jwtS: jwtS,
    }
}

func (m *JWTAuth) Handler(guardName string) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenStr := c.Request.Header.Get("Authorization")
        if tokenStr == "" {
            response.FailByErr(c, cErr.Unauthorized("missing Authorization header"))
            return
        }
        tokenStr = tokenStr[len(domain.TokenType)+1:]

        token, err := jwt.ParseWithClaims(tokenStr, &domain.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
            return []byte(m.conf.Jwt.Secret), nil
        })
        if err != nil || m.jwtS.IsInBlacklist(c, tokenStr) {
            response.FailByErr(c, cErr.Unauthorized("登录授权已失效"))
            return
        }

        claims := token.Claims.(*domain.CustomClaims)
        if claims.Issuer != guardName {
            response.FailByErr(c, cErr.Unauthorized("登录授权已失效"))
            return
        }

        // token 续签
        if claims.ExpiresAt-time.Now().Unix() < m.conf.Jwt.RefreshGracePeriod {
            tokenData, err := m.jwtS.RefreshToken(c, guardName, token)
            if err == nil {
                c.Header("new-token", tokenData.AccessToken)
                c.Header("new-expires-in", strconv.Itoa(tokenData.ExpiresIn))
            }
        }

        c.Set("token", token)
        c.Set("id", claims.Id)
    }
}