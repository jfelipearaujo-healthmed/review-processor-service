package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	event_processor_uc "github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/application/use_cases/event/event_processor"
	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/infrastructure/config"
	appointment_repository "github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/infrastructure/repository/appointment"
	event_repository "github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/infrastructure/repository/event"
	feedback_repository "github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/infrastructure/repository/feedback"
	user_repository "github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/infrastructure/repository/user"
	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/external/http/handlers/health"
	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/external/http/handlers/middlewares/logger"
	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/external/persistence"
	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/external/queue"
	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/external/secret"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Config *config.Config

	Dependencies
}

func NewServer(ctx context.Context, config *config.Config) (*Server, error) {
	cloudConfig, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "error getting aws config", "error", err)
		return nil, err
	}

	if config.CloudConfig.IsBaseEndpointSet() {
		cloudConfig.BaseEndpoint = aws.String(config.CloudConfig.BaseEndpoint)
	}

	secretService := secret.NewService(cloudConfig)

	dbUrl, err := secretService.GetSecret(ctx, config.DbConfig.UrlSecretName)
	if err != nil {
		slog.ErrorContext(ctx, "error getting secret", "secret_name", config.DbConfig.UrlSecretName, "error", err)
		return nil, err
	}

	config.DbConfig.Url = dbUrl

	dbService := persistence.NewDbService()

	if err := dbService.Connect(config); err != nil {
		slog.ErrorContext(ctx, "error connecting to database", "error", err)
		return nil, err
	}

	eventRepository := event_repository.NewRepository(dbService)
	appointmentRepository := appointment_repository.NewRepository(dbService)
	feedbackRepository := feedback_repository.NewRepository(dbService)
	userRepository := user_repository.NewRepository(config)

	eventProcessor := event_processor_uc.NewUseCase(eventRepository, appointmentRepository, feedbackRepository, userRepository)

	reviewQueueService := queue.NewService(config.CloudConfig.ReviewQueueName, cloudConfig, eventProcessor.Handle)

	if err := reviewQueueService.UpdateQueueUrl(ctx); err != nil {
		slog.ErrorContext(ctx, "error updating queue url", "error", err)
		return nil, err
	}

	return &Server{
		Config: config,
		Dependencies: Dependencies{
			DbService: dbService,

			ReviewQueueService: reviewQueueService,
		},
	}, nil
}

func (s *Server) GetServer() *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", s.Config.ApiConfig.Port),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(logger.Middleware())
	e.Use(middleware.Recover())

	s.addHealthCheckRoutes(e)

	return e
}

func (s *Server) addHealthCheckRoutes(e *echo.Echo) {
	healthHandler := health.NewHandler(s.DbService)

	e.GET("/health", healthHandler.Handle)
}
