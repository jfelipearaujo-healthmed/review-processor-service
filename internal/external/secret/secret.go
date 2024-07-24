package secret

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type Service struct {
	client *secretsmanager.Client
}

func NewService(config aws.Config) *Service {
	return &Service{
		client: secretsmanager.NewFromConfig(config),
	}
}

func (s *Service) GetSecret(ctx context.Context, secretName string) (string, error) {
	output, err := s.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"),
	})

	if err != nil {
		return "", err
	}

	return *output.SecretString, nil
}
