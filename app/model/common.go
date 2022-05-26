package model

import (
    "gorm.io/gorm"
    "time"
)

type U64ID struct {
    ID uint64 `gorm:"primaryKey"`
}

type Timestamps struct {
    CreatedAt time.Time `gorm:"column:created_at"`
    UpdatedAt time.Time `gorm:"column:updated_at"`
}

type SoftDeletes struct {
    DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at"`
}
