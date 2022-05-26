package common

import (
    "github.com/gin-gonic/gin"
    "github.com/jassue/gin-wire/app/pkg/request"
    "github.com/jassue/gin-wire/app/pkg/response"
    "github.com/jassue/gin-wire/app/service"
    "go.uber.org/zap"
)

type UploadHandler struct {
    log *zap.Logger
    mediaS *service.MediaService
}

func NewUploadHandler(log *zap.Logger, mediaS *service.MediaService) *UploadHandler {
    return &UploadHandler{log: log, mediaS: mediaS}
}

func (h *UploadHandler) ImageUpload(c *gin.Context) {
    var form request.ImageUpload
    if err := c.ShouldBind(&form); err != nil {
        response.FailByErr(c, request.GetError(form, err))
        return
    }

    media, err := h.mediaS.SaveImage(c, &form)
    if err != nil {
        response.FailByErr(c, err)
        return
    }

    response.Success(c, media)
}
