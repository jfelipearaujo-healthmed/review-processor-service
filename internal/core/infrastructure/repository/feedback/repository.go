package feedback_repository

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/domain/entities"
	feedback_repository_contract "github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/domain/repositories/feedback"
	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/infrastructure/shared/app_error"
	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/external/persistence"
	"gorm.io/gorm"
)

type repository struct {
	dbService *persistence.DbService
}

func NewRepository(dbService *persistence.DbService) feedback_repository_contract.Repository {
	return &repository{
		dbService: dbService,
	}
}

func (rp *repository) GetByID(ctx context.Context, appointmentID, feedbackID uint) (*entities.Feedback, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	feedback := new(entities.Feedback)

	if err := tx.Preload("Appointment").
		Order("feedbacks.created_at DESC").
		Joins("JOIN appointments ON appointments.id = feedbacks.appointment_id").
		Where("feedbacks.id = ? AND appointments.id = ?", feedbackID, appointmentID).
		Find(&feedback).Error; err != nil {
		return nil, err
	}

	return feedback, nil
}

func (rp *repository) GetByAppointmentID(ctx context.Context, appointmentID uint) (*entities.Feedback, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	feedback := new(entities.Feedback)

	if err := tx.Where("appointment_id = ?", appointmentID).First(feedback).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_error.New(http.StatusNotFound, fmt.Sprintf("feedback with appointment id %d not found", appointmentID))
		}

		return nil, err
	}

	return feedback, nil
}

func (rp *repository) GetRatingAndCountFromDoctor(ctx context.Context, doctorID uint) (float64, int, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	var result struct {
		AvgRating float64 `json:"avg_rating"`
		Count     int     `json:"count"`
	}

	if err := tx.Model(&entities.Feedback{}).
		Select("avg(feedbacks.rating) as avg_rating, count(distinct feedbacks.id) as count").
		Joins("JOIN appointments ON feedbacks.appointment_id = appointments.id").
		Where("appointments.doctor_id = ?", doctorID).
		Find(&result).Error; err != nil {
		return 0, 0, err
	}

	return result.AvgRating, result.Count, nil
}

func (rp *repository) Create(ctx context.Context, feedback *entities.Feedback) (*entities.Feedback, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Create(feedback).Error; err != nil {
		return nil, err
	}

	return feedback, nil
}

func (rp *repository) Update(ctx context.Context, feedback *entities.Feedback) (*entities.Feedback, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Model(feedback).Where("id = ?", feedback.ID).Updates(feedback).Error; err != nil {
		return nil, err
	}

	return feedback, nil
}

func (rp *repository) Delete(ctx context.Context, feedbackID uint) error {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Delete(&entities.Feedback{}, "id = ?", feedbackID).Error; err != nil {
		return err
	}

	return nil
}
