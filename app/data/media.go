package data

import (
    "context"
    "encoding/json"
    "github.com/jassue/gin-wire/app/compo"
    "github.com/jassue/gin-wire/app/domain"
    "github.com/jassue/gin-wire/app/model"
    "github.com/jassue/gin-wire/app/service"
    "go.uber.org/zap"
    "strconv"
    "time"
)

const mediaCacheKeyPre = "media:"

type mediaRepo struct {
    data *Data
    log *zap.Logger
    storage *compo.Storage
}

func NewMediaRepo(data *Data, log *zap.Logger, storage *compo.Storage) service.MediaRepo {
    return &mediaRepo{
        data: data,
        log: log,
        storage: storage,
    }
}

func (r *mediaRepo) Create(ctx context.Context, dm *domain.Media) (*domain.Media, error) {
    var m model.Media

    id, err := r.data.sf.NextID()
    if err != nil {
        return nil, err
    }

    m.ID = id
    m.DiskType = dm.DiskType
    m.SrcType = dm.SrcType
    m.Src = dm.Src

    if err = r.data.DB(ctx).Create(&m).Error; err != nil {
        return nil, err
    }
    dm.ID = m.ID

    return dm, nil
}

func (r *mediaRepo) FindByID(ctx context.Context, id uint64) (*domain.Media, error) {
    var m model.Media
    if err := r.data.db.First(&m, id).Error; err != nil{
        return nil, err
    }

    dm := m.ToDomain()
    dm.SetUrl(r.storage)

    return dm, nil
}

func (r *mediaRepo) FindCacheByID(ctx context.Context, id uint64) (*domain.Media, error) {
    cacheKey := mediaCacheKeyPre + strconv.FormatUint(id,10)

    exist := r.data.rdb.Exists(ctx, cacheKey).Val()
    if exist == 1 {
        bytes, err := r.data.rdb.Get(ctx, cacheKey).Bytes()
        if err != nil {
            return nil, err
        }
        var media domain.Media
        err = json.Unmarshal(bytes, &media)
        if err != nil {
            return nil, err
        }

        return &media, nil
    }

    var media model.Media
    err := r.data.db.First(&media, id).Error
    if err != nil {
        return nil, err
    }
    dm := media.ToDomain()
    dm.SetUrl(r.storage)
    v, err := json.Marshal(dm)
    if err != nil {
        return nil, err
    }
    r.data.rdb.Set(ctx, cacheKey, v, time.Second*3*24*3600)

    return dm, nil
}

