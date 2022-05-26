package utils

import (
    "os"
    "path/filepath"
    "runtime"
    "strings"
)

// RootPath 获取项目根目录绝对路径
func RootPath() string {
    var rootDir string

    exePath, err := os.Executable()
    if err != nil {
        panic(err)
    }

    rootDir = filepath.Dir(filepath.Dir(exePath))

    tmpDir := os.TempDir()
    if strings.Contains(exePath, tmpDir) {
        _, filename, _, ok := runtime.Caller(0)
        if ok {
            rootDir = filepath.Dir(filepath.Dir(filename))
        }
    }

    return rootDir
}
