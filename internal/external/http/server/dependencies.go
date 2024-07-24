package server

import (
	appointment_repository_contract "github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/domain/repositories/appointment"
	event_repository_contract "github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/domain/repositories/event"
	feedback_repository_contract "github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/domain/repositories/feedback"
	event_processor_contract "github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/domain/use_cases/event/event_processor"
	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/external/persistence"
	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/external/queue"
)

type Dependencies struct {
	DbService *persistence.DbService

	ReviewQueueService queue.QueueService

	EventRepository       event_repository_contract.Repository
	AppointmentRepository appointment_repository_contract.Repository
	FeedbackRepository    feedback_repository_contract.Repository

	EventProcessor event_processor_contract.UseCase
}
