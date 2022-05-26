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

func JWTAuth(conf *config.Jwt, jwtS *service.JwtService, guardName string) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenStr := c.Request.Header.Get("Authorization")
        if tokenStr == "" {
            response.FailByErr(c, cErr.Unauthorized("missing Authorization header"))
            c.Abort()
            return
        }
        tokenStr = tokenStr[len(domain.TokenType)+1:]

        token, err := jwt.ParseWithClaims(tokenStr, &domain.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
            return []byte(conf.Secret), nil
        })
        if err != nil || jwtS.IsInBlacklist(c, tokenStr) {
            response.FailByErr(c, cErr.Unauthorized("登录授权已失效"))
            c.Abort()
            return
        }

        claims := token.Claims.(*domain.CustomClaims)
        if claims.Issuer != guardName {
            response.FailByErr(c, cErr.Unauthorized("登录授权已失效"))
            c.Abort()
            return
        }

        // token 续签
        if claims.ExpiresAt-time.Now().Unix() < conf.RefreshGracePeriod {
            tokenData, err := jwtS.RefreshToken(c, guardName, token)
            if err == nil {
                c.Header("new-token", tokenData.AccessToken)
                c.Header("new-expires-in", strconv.Itoa(tokenData.ExpiresIn))
            }
        }

        c.Set("token", token)
        c.Set("id", claims.Id)
    }
}