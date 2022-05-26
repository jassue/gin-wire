package command

import (
    "github.com/jassue/gin-wire/app/model"
    "github.com/spf13/cobra"
    "go.uber.org/zap"
    "gorm.io/gorm"
)

type MigrateHandler struct {
    logger *zap.Logger
    db *gorm.DB
}

func NewMigrateHandler(logger *zap.Logger, db *gorm.DB) *MigrateHandler {
    return &MigrateHandler{
        logger: logger,
        db: db,
    }
}

func (h *MigrateHandler) Migrate(cmd *cobra.Command, args []string) {
    err := h.db.AutoMigrate(
        &model.User{},
        &model.Media{},
        )

    if err != nil {
        cmd.Println("database migrate error:", err)
    }

    cmd.Println("database migrate success")
}
