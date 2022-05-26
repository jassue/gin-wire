package model

import (
    "gorm.io/gorm"
    "time"
)

type U64ID struct {
    ID uint64 `json:"id" gorm:"primaryKey"`
}

type Timestamps struct {
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type SoftDeletes struct {
    DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
