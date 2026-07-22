package jobs

import (
	logging "github.com/gcc798/lightning/internal/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Definitions returns the code-owned jobs that a scheduler process may run.
func Definitions(
	db *gorm.DB,
	logger logging.Logger,
) map[string]func() {
	definitions := make(map[string]func())
	cl := NewDataCleanupJob(db, logger)
	definitions["data-cleanup"] = cl.Run

	logger.Info("scheduler job definitions created", zap.Int("count", len(definitions)))
	return definitions
}
