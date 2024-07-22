package zap_info_logger

import (
	"time"

	"go.uber.org/zap"
)

func zap_info_logger() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger = logger.Named("my-app-logs")
	logger.Info(
		"failed to fetch the URL",
		zap.String("url", "https://github.com/lyteabovenyte"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
}
