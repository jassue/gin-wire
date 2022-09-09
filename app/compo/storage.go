package compo

import (
    "github.com/jassue/gin-wire/config"
    "github.com/jassue/gin-wire/util/path"
    "github.com/jassue/go-storage/kodo"
    "github.com/jassue/go-storage/local"
    "github.com/jassue/go-storage/oss"
    "github.com/jassue/go-storage/storage"
    "go.uber.org/zap"
    "path/filepath"
)

type Storage struct {
    conf *config.Configuration
    log *zap.Logger
}

// NewStorage .
func NewStorage(c *config.Configuration, log *zap.Logger) *Storage {
    if len(c.Storage.Default) == 0 {
        panic("disk config error")
    }

    if !filepath.IsAbs(c.Storage.Disks.Local.RootDir) {
        c.Storage.Disks.Local.RootDir = filepath.Join(path.RootPath(), c.Storage.Disks.Local.RootDir)
    }
    _, _ = local.Init(c.Storage.Disks.Local)
    _, _ = kodo.Init(c.Storage.Disks.QiNiu)
    _, _ = oss.Init(c.Storage.Disks.AliOss)

    return &Storage{conf: c, log: log}
}

func (s *Storage) IsLocal() bool {
    return s.conf.Storage.Default == storage.Local
}

func (s *Storage) GetDefaultDiskType() string {
    return string(s.conf.Storage.Default)
}

func (s *Storage) GetDisk(disk... string) (storage.Storage, error) {
    diskName := s.conf.Storage.Default
    if len(disk) > 0 {
        diskName = storage.DiskName(disk[0])
    }
    d, err := storage.Disk(diskName)
    if err != nil {
        s.log.Error(err.Error())
        return nil, err
    }

    return d, nil
}
