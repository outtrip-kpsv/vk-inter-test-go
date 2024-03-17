package bl

import (
	"go.uber.org/zap"
	"vk-inter-test-go/internal/db"
)

type BL struct {
	Db     *db.DBRepo
	logger *zap.Logger
}

func NewBL(repo *db.DBRepo, logger *zap.Logger) *BL {
	logger = logger.Named("Bl")

	return &BL{
		Db:     repo,
		logger: logger,
	}
}
