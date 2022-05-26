package service

import (
    "context"
    "github.com/gin-gonic/gin"
    "github.com/jassue/gin-wire/app/compo"
    "github.com/jassue/gin-wire/app/domain"
    cErr "github.com/jassue/gin-wire/app/pkg/error"
    "github.com/jassue/gin-wire/app/pkg/request"
    "github.com/jassue/gin-wire/config"
    "github.com/satori/go.uuid"
    "go.uber.org/zap"
    "path"
)

type MediaRepo interface {
   Create(context.Context, *domain.Media) (*domain.Media, error)
   FindByID(context.Context, uint64) (*domain.Media, error)
   FindCacheByID(context.Context, uint64) (*domain.Media, error)
}

type MediaService struct {
    conf *config.Configuration
    log *zap.Logger
    mRepo MediaRepo
    storage *compo.Storage
}

// NewMediaService .
func NewMediaService(conf *config.Configuration, log *zap.Logger, mRepo MediaRepo, s *compo.Storage) *MediaService {
    return &MediaService{conf: conf, log: log, mRepo: mRepo, storage: s}
}

func (s *MediaService) makeFaceDir(business string) string {
    return s.conf.App.Env + "/" + business
}

func (s *MediaService) HashName(fileName string) string {
    fileSuffix := path.Ext(fileName)
    return uuid.NewV4().String() + fileSuffix
}

// SaveImage 保存图片（公共读）
func (s *MediaService) SaveImage(ctx *gin.Context, params *request.ImageUpload) (*domain.Media, error) {
    file, err := params.Image.Open()
    defer file.Close()
    if err != nil {
        return nil, cErr.BadRequest("上传失败")
    }

    disk, err := s.storage.GetDisk()
    if err != nil {
        return nil, cErr.BadRequest(s.storage.GetDefaultDiskType() + "disk not found")
    }
    localPrefix := ""
    if s.storage.IsLocal() {
        localPrefix = "public" + "/"
    }
    key := s.makeFaceDir(params.Business) + "/" + s.HashName(params.Image.Filename)
    err = disk.Put(localPrefix + key, file, params.Image.Size)
    if err != nil {
        return nil, cErr.BadRequest("上传失败")
    }

    m, err := s.mRepo.Create(ctx, &domain.Media{
        DiskType: s.storage.GetDefaultDiskType(),
        SrcType:  1,
        Src:      key,
        Url:      disk.Url(key),
    })
    if err != nil {
        return nil, cErr.BadRequest("上传失败")
    }

    return m, nil
}

func (s *MediaService) GetUrlById(ctx *gin.Context, id uint64) string {
    if id == 0 {
        return ""
    }
    m, err := s.mRepo.FindCacheByID(ctx, id)
    if err != nil {
        s.log.Error(err.Error())
        return ""
    }

    disk, err := s.storage.GetDisk(m.DiskType)
    if err != nil {
        s.log.Error(err.Error())
        return ""
    }

    return disk.Url(m.Src)
}
