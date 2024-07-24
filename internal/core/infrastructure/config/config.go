package config

import (
	"context"
	"time"

	"github.com/sethvargo/go-envconfig"
)

type ApiConfig struct {
	Port           int    `env:"PORT, default=5000"`
	EnvName        string `env:"ENV_NAME, default=development"`
	ApiVersion     string `env:"VERSION, default=v1"`
	LocationRegion string `env:"LOCATION, default=America/Sao_Paulo"`
	Location       *time.Location
}

func (c *ApiConfig) IsDevelopment() bool {
	return c.EnvName == "development"
}

type DatabaseConfig struct {
	Url           string `env:"URL, required"`
	UrlSecretName string `env:"URL_SECRET_NAME, required"`
}

type CloudConfig struct {
	BaseEndpoint    string `env:"BASE_ENDPOINT"`
	ReviewQueueName string `env:"REVIEW_QUEUE_NAME"`
}

type ExternalServiceConfig struct {
	UserServiceBaseRoute string `env:"USER_SERVICE_BASE_ROUTE"`
}

func (c *CloudConfig) IsBaseEndpointSet() bool {
	return c.BaseEndpoint != ""
}

type Config struct {
	ApiConfig             *ApiConfig             `env:",prefix=API_"`
	DbConfig              *DatabaseConfig        `env:",prefix=DB_"`
	CloudConfig           *CloudConfig           `env:",prefix=AWS_"`
	ExternalServiceConfig *ExternalServiceConfig `env:",prefix=EXT_"`
}

func LoadFromEnv(ctx context.Context) (*Config, error) {
	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		return nil, err
	}

	Location, err := time.LoadLocation(cfg.ApiConfig.LocationRegion)
	if err != nil {
		return nil, err
	}

	cfg.ApiConfig.Location = Location

	return &cfg, nil
}
