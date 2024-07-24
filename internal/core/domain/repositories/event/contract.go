package event_repository_contract

import (
	"context"

	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/domain/entities"
)

type Repository interface {
	GetByMessageID(ctx context.Context, messageID string) (*entities.Event, error)
	GetByIDsAndDateTime(ctx context.Context, event *entities.Event) (*entities.Event, error)
	Create(ctx context.Context, event *entities.Event) (*entities.Event, error)
	Update(ctx context.Context, event *entities.Event) (*entities.Event, error)
}
