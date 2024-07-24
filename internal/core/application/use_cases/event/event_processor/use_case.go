package event_processor_uc

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/domain/entities"
	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/domain/events"
	appointment_repository_contract "github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/domain/repositories/appointment"
	event_repository_contract "github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/domain/repositories/event"
	feedback_repository_contract "github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/domain/repositories/feedback"
	user_repository_contract "github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/domain/repositories/user"
	event_processor_contract "github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/domain/use_cases/event/event_processor"
	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/external/queue"
)

var (
	OUTCOME_FEEDBACK_PROCESSED string = "feedback processed successfully"
)

type useCase struct {
	eventRepository       event_repository_contract.Repository
	appointmentRepository appointment_repository_contract.Repository
	feedbackRepository    feedback_repository_contract.Repository
	userRepository        user_repository_contract.Repository
}

func NewUseCase(
	eventRepository event_repository_contract.Repository,
	appointmentRepository appointment_repository_contract.Repository,
	feedbackRepository feedback_repository_contract.Repository,
	userRepository user_repository_contract.Repository,
) event_processor_contract.UseCase {
	return &useCase{
		eventRepository:       eventRepository,
		appointmentRepository: appointmentRepository,
		feedbackRepository:    feedbackRepository,
		userRepository:        userRepository,
	}
}

func (uc *useCase) Handle(ctx context.Context, messageID string, message queue.Message) error {
	slog.InfoContext(ctx, "event received", "message_id", messageID)

	event := new(entities.Event)

	messageJson, err := json.Marshal(message.Data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(messageJson, event); err != nil {
		return err
	}

	eventMap := map[queue.EventType]func(ctx context.Context, messageID string, feedback *entities.Feedback) error{
		events.CreateFeedback: uc.HandleCreateReview,
	}

	slog.InfoContext(ctx, "checking handler for event type", "event_type", event.EventType)

	if handler, ok := eventMap[event.EventType]; ok {
		feedback := new(entities.Feedback)

		if err := json.Unmarshal([]byte(event.Data), feedback); err != nil {
			return err
		}

		err = handler(ctx, messageID, feedback)
		if err != nil {
			return err
		}

		slog.InfoContext(ctx, "event processed successfully", "message_id", messageID)

		return nil
	}

	slog.ErrorContext(ctx, "event handler not found", "message_id", messageID)

	return nil
}

func (uc *useCase) HandleCreateReview(ctx context.Context, messageID string, feedback *entities.Feedback) error {
	slog.InfoContext(ctx, "handling feedback creation", "message_id", messageID)

	slog.InfoContext(ctx, "loading event for message received", "message_id", messageID)

	event, err := uc.eventRepository.GetByMessageID(ctx, messageID)
	if err != nil {
		slog.ErrorContext(ctx, "error loading event", "message_id", messageID, "error", err)
		return err
	}

	if event.Outcome != nil {
		slog.WarnContext(ctx, "event already processed", "message_id", messageID)
		return nil
	}

	slog.InfoContext(ctx, "loading appointment for feedback", "message_id", messageID)

	appointment, err := uc.appointmentRepository.GetByID(ctx, feedback.AppointmentID)
	if err != nil {
		slog.ErrorContext(ctx, "error checking if appointment already exists", "message_id", messageID, "error", err)
		return err
	}

	slog.InfoContext(ctx, "loading doctor for feedback", "message_id", messageID, "doctor_id", appointment.DoctorID)

	doctor, err := uc.userRepository.GetByDoctorID(ctx, appointment.DoctorID)
	if err != nil {
		slog.ErrorContext(ctx, "error loading doctor", "message_id", messageID, "error", err)
		return err
	}

	currentAverage, currentTotal, err := uc.feedbackRepository.GetRatingAndCountFromDoctor(ctx, doctor.ID)
	if err != nil {
		slog.ErrorContext(ctx, "error loading doctor", "message_id", messageID, "error", err)
		return err
	}

	newTotal := currentTotal + 1
	newAverage := ((currentAverage * float64(currentTotal)) + feedback.Rating) / float64(newTotal)

	slog.InfoContext(ctx, "updating doctor rating", "message_id", messageID, "doctor_id", doctor.ID, "rating", newAverage)

	if err := uc.userRepository.UpdateRating(ctx, doctor.ID, newAverage); err != nil {
		slog.ErrorContext(ctx, "error updating doctor rating", "message_id", messageID, "error", err)
		return err
	}

	slog.InfoContext(ctx, "creating feedback", "message_id", messageID)

	if _, err := uc.feedbackRepository.Create(ctx, feedback); err != nil {
		slog.ErrorContext(ctx, "error creating feedback", "message_id", messageID, "error", err)
		return err
	}

	event.Outcome = &OUTCOME_FEEDBACK_PROCESSED

	return uc.updateEvent(ctx, messageID, event)
}

func (uc *useCase) updateEvent(ctx context.Context, messageID string, event *entities.Event) error {
	slog.InfoContext(ctx, "updating event...", "message_id", messageID)

	if _, err := uc.eventRepository.Update(ctx, event); err != nil {
		slog.ErrorContext(ctx, "error updating event", "message_id", messageID, "error", err)
		return err
	}

	slog.InfoContext(ctx, "event updated successfully", "message_id", messageID)

	return nil
}
