package domain

import "github.com/golang-jwt/jwt"

const (
    TokenType = "Bearer"
    AppGuardName = "app"
)

type JwtUser interface {
    GetUid() string
}

// CustomClaims 自定义 Claims
type CustomClaims struct {
    jwt.StandardClaims
}

type TokenOutPut struct {
    AccessToken string `json:"access_token"`
    ExpiresIn int `json:"expires_in"`
    TokenType string `json:"token_type"`
}
