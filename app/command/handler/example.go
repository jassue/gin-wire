package command

import (
    "github.com/spf13/cobra"
    "go.uber.org/zap"
    "strings"
)

type ExampleHandler struct {
    logger *zap.Logger
}

func NewExampleHandler(logger *zap.Logger) *ExampleHandler {
    return &ExampleHandler{
        logger: logger,
    }
}

func (h *ExampleHandler) Hello(cmd *cobra.Command, args []string) {
    cmd.Println(cmd.Use, "命令调用成功")
    cmd.Printf("Hello %s\n", strings.Join(args, ","))
}