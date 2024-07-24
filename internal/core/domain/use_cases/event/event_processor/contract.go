package event_processor_contract

import (
	"context"

	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/external/queue"
)

type UseCase interface {
	Handle(ctx context.Context, messageID string, message queue.Message) error
}
