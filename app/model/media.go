package model

import (
    "github.com/jassue/gin-wire/app/domain"
)

type Media struct {
    U64ID
    DiskType string `gorm:"size:20;index;not null;comment:存储类型"`
    SrcType int8 `gorm:"not null;comment:链接类型 1相对路径 2外链"`
    Src string `gorm:"not null;comment:资源链接"`
    Timestamps
}

func (m *Media) ToDomain() *domain.Media {
    return &domain.Media{
        ID:       m.ID,
        DiskType: m.DiskType,
        SrcType:  m.SrcType,
        Src:      m.Src,
    }
}
