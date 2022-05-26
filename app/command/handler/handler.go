package command

import (
    "github.com/google/wire"
    "github.com/jassue/gin-wire/app/data"
)

// ProviderSet is handler providers.
var ProviderSet = wire.NewSet(data.NewDB, NewExampleHandler, NewMigrateHandler)
