package compo

import (
    "github.com/google/wire"
    "github.com/jassue/gin-wire/app/compo/casbin"
)

// ProviderSet is compo providers.
var ProviderSet = wire.NewSet(
    NewSonyFlake,
    NewLockBuilder,
    NewStorage,
    casbin.NewEnforcer,
    )
