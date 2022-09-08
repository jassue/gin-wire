package casbin

import (
    "github.com/casbin/casbin/util"
    "github.com/casbin/casbin/v2"
    gormadapter "github.com/casbin/gorm-adapter/v3"
    "github.com/jassue/gin-wire/utils/path"
    "gorm.io/gorm"
    "path/filepath"
    "strings"
)

func NewEnforcer(db *gorm.DB) *casbin.Enforcer {
    adapter, err := gormadapter.NewAdapterByDB(db)
    if err != nil {
        panic(err)
    }

    enforcer, err := casbin.NewEnforcer(filepath.Join(path.RootPath(), "app/compo/casbin/rbac_model.conf"), adapter)
    if err != nil {
        panic(err)
    }

    //enforcer.AddFunction("ParamsMatch", ParamsMatchFunc)

    _ = enforcer.LoadPolicy()
    return enforcer
}

func ParamsMatch(fullNameKey1 string, key2 string) bool {
    key1 := strings.Split(fullNameKey1, "?")[0]
    // 剥离路径后再使用casbin的keyMatch2
    return util.KeyMatch2(key1, key2)
}

func ParamsMatchFunc(args ...interface{}) (interface{}, error) {
    name1 := args[0].(string)
    name2 := args[1].(string)

    return ParamsMatch(name1, name2), nil
}