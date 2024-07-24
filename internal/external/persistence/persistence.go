package persistence

import (
	"context"
	"log/slog"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/infrastructure/config"
	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/infrastructure/shared/logger"
	slogGorm "github.com/orandin/slog-gorm"
)

type DbService struct {
	Instance *gorm.DB
}

func NewDbService() *DbService {
	return &DbService{}
}

func (svc *DbService) Connect(config *config.Config) error {
	log := slog.Default()

	gormLogger := slogGorm.New(
		slogGorm.WithHandler(log.Handler()),
		slogGorm.WithIgnoreTrace(),
		slogGorm.WithSlowThreshold(500*time.Millisecond),
		slogGorm.SetLogLevel(slogGorm.ErrorLogType, logger.LOG_ERR),
		slogGorm.SetLogLevel(slogGorm.SlowQueryLogType, logger.LOG_NOTICE),
		slogGorm.SetLogLevel(slogGorm.DefaultLogType, logger.LOG_DEBUG),
	)

	conn, err := gorm.Open(postgres.Open(config.DbConfig.Url), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return err
	}

	svc.Instance = conn

	return nil
}

func (svc *DbService) Close(ctx context.Context) error {
	slog.InfoContext(ctx, "closing database connection")

	db, err := svc.Instance.DB()
	if err != nil {
		slog.ErrorContext(ctx, "error closing database connection", "error", err)
		return err
	}

	if err := db.Close(); err != nil {
		slog.ErrorContext(ctx, "error closing database connection", "error", err)
		return err
	}

	slog.InfoContext(ctx, "database connection closed")

	return nil
}
