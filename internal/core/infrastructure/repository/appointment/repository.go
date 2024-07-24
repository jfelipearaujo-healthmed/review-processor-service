package appointment_repository

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/domain/entities"
	appointment_repository_contract "github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/domain/repositories/appointment"
	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/infrastructure/shared/app_error"
	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/external/persistence"
	"gorm.io/gorm"
)

type repository struct {
	dbService *persistence.DbService
}

func NewRepository(dbService *persistence.DbService) appointment_repository_contract.Repository {
	return &repository{
		dbService: dbService,
	}
}

func (rp *repository) GetByID(ctx context.Context, appointmentID uint) (*entities.Appointment, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	appointment := new(entities.Appointment)

	if err := tx.Where("id = ?", appointmentID).First(appointment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app_error.New(http.StatusNotFound, fmt.Sprintf("appointment with id %d not found", appointmentID))
		}

		return nil, err
	}

	return appointment, nil
}

func (rp *repository) GetByIDsAndDateTime(ctx context.Context, scheduleID uint, doctorID uint, dateTime time.Time) (*entities.Appointment, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	appointment := new(entities.Appointment)

	query := tx.Where("schedule_id = ?", scheduleID)
	query = query.Where("doctor_id = ?", doctorID)
	query = query.Where("date_time = ?", dateTime)
	query = query.Where("status != ?", entities.Cancelled)

	result := query.First(appointment)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, app_error.New(http.StatusNotFound, fmt.Sprintf("appointment with schedule id %d, doctor id %d and date time %s not found", scheduleID, doctorID, dateTime))
		}

		return nil, result.Error
	}

	return appointment, nil
}

func (rp *repository) Create(ctx context.Context, appointment *entities.Appointment) (*entities.Appointment, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Create(appointment).Error; err != nil {
		return nil, err
	}

	return appointment, nil
}

func (rp *repository) Update(ctx context.Context, userID uint, appointment *entities.Appointment) (*entities.Appointment, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Model(appointment).Where("patient_id = ? AND id = ?", userID, appointment.ID).Updates(appointment).Error; err != nil {
		return nil, err
	}

	return appointment, nil
}

func (rp *repository) Delete(ctx context.Context, userID uint, appointmentID uint) error {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Delete(&entities.Appointment{}, "patient_id = ? AND id = ?", userID, appointmentID).Error; err != nil {
		return err
	}

	return nil
}
