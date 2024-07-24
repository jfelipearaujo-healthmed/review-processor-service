package feedback_repository_contract

import (
	"context"

	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/domain/entities"
)

type Repository interface {
	GetByID(ctx context.Context, appointmentID, feedbackID uint) (*entities.Feedback, error)
	GetByAppointmentID(ctx context.Context, appointmentID uint) (*entities.Feedback, error)
	GetRatingAndCountFromDoctor(ctx context.Context, doctorID uint) (float64, int, error)
	Create(ctx context.Context, feedback *entities.Feedback) (*entities.Feedback, error)
	Update(ctx context.Context, feedback *entities.Feedback) (*entities.Feedback, error)
	Delete(ctx context.Context, feedbackID uint) error
}
