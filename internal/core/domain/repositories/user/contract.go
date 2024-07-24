package user_repository_contract

import (
	"context"

	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/domain/entities"
)

type Repository interface {
	GetByDoctorID(ctx context.Context, patientID, doctorID uint) (*entities.Doctor, error)
	UpdateRating(ctx context.Context, patientID, doctorID uint, rating float64) error
}
