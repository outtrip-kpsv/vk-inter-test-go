package handlers

import (
	"go.uber.org/zap"
	"vk-inter-test-go/internal/bl"
)

type Controller struct {
	Bl     *bl.BL
	logger *zap.Logger
}

func NewController(bl *bl.BL, logger *zap.Logger) *Controller {
	logger = logger.Named("Handler")
	return &Controller{Bl: bl, logger: logger}
}
