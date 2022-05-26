package handler

import (
    "github.com/google/wire"
    "github.com/jassue/gin-wire/app/handler/app"
    "github.com/jassue/gin-wire/app/handler/common"
)

// ProviderSet is handler providers.
var ProviderSet = wire.NewSet(
    app.NewAuthHandler,
    common.NewUploadHandler,
    )
