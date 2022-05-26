package router

import (
    "github.com/gin-gonic/gin"
    "github.com/jassue/gin-wire/app/domain"
    "github.com/jassue/gin-wire/app/handler/app"
    "github.com/jassue/gin-wire/app/handler/common"
    "github.com/jassue/gin-wire/app/middleware"
    "github.com/jassue/gin-wire/app/service"
    "github.com/jassue/gin-wire/config"
)

func setApiGroupRoutes(
    router *gin.Engine,
    conf *config.Configuration,
    jwtS *service.JwtService,
    authH *app.AuthHandler,
    commonH *common.UploadHandler,
    ) *gin.RouterGroup {
    group := router.Group("/api")
    group.POST("/auth/register", authH.Register)
    group.POST("/auth/login", authH.Login)
    authGroup := group.Group("").Use(middleware.JWTAuth(&conf.Jwt, jwtS, domain.AppGuardName))
    {
        authGroup.POST("/auth/info", authH.Info)
        authGroup.POST("/auth/logout", authH.Logout)
        authGroup.POST("/image_upload", commonH.ImageUpload)
    }

    return group
}
