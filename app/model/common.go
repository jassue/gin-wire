package model

import (
    "gorm.io/gorm"
    "time"
)

type Timestamps struct {
    CreatedAt time.Time `gorm:"column:created_at"`
    UpdatedAt time.Time `gorm:"column:updated_at"`
}

type SoftDeletes struct {
    DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at"`
}
