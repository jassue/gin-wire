package app

import (
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt"
    "github.com/jassue/gin-wire/app/domain"
    "github.com/jassue/gin-wire/app/pkg/request"
    "github.com/jassue/gin-wire/app/pkg/response"
    "github.com/jassue/gin-wire/app/service"
    "go.uber.org/zap"
)

type AuthHandler struct {
    log *zap.Logger
    jwtS *service.JwtService
    userS *service.UserService
}

func NewAuthHandler(log *zap.Logger, jwtS *service.JwtService, userS *service.UserService) *AuthHandler {
    return &AuthHandler{log: log, jwtS: jwtS, userS: userS}
}

func (h *AuthHandler) Register(c *gin.Context) {
    var form request.Register
    if err := c.ShouldBindJSON(&form); err != nil {
        response.FailByErr(c, request.GetError(form, err))
        return
    }

    u, err := h.userS.Register(c, &form)
    if err != nil {
        response.FailByErr(c, err)
        return
    }

    tokenData, _, err := h.jwtS.CreateToken(domain.AppGuardName, u)
    if err != nil {
        response.FailByErr(c, err)
        return
    }

    response.Success(c, tokenData)
}

func (h *AuthHandler) Login(c *gin.Context) {
    var form request.Login
    if err := c.ShouldBindJSON(&form); err != nil {
        response.FailByErr(c, request.GetError(form, err))
        return
    }

    user, err := h.userS.Login(c, form.Mobile, form.Password)
    if err != nil {
        response.FailByErr(c, err)
        return
    }

    tokenData, _, err := h.jwtS.CreateToken(domain.AppGuardName, user)
    if err != nil {
        response.FailByErr(c, err)
        return
    }

    response.Success(c, tokenData)
}

func (h *AuthHandler) Info(c *gin.Context) {
    user, err := h.userS.GetUserInfo(c, c.Keys["id"].(string))
    if err != nil {
        response.FailByErr(c, err)
        return
    }
    response.Success(c, user)
}

func (h *AuthHandler) Logout(c *gin.Context) {
    err := h.jwtS.JoinBlackList(c, c.Keys["token"].(*jwt.Token))
    if err != nil {
        response.FailByErr(c, err)
        return
    }

    response.Success(c, nil)
}
