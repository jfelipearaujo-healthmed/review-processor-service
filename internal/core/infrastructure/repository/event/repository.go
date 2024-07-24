package event_repository

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/domain/entities"
	event_repository_contract "github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/domain/repositories/event"
	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/infrastructure/shared/app_error"
	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/external/persistence"
	"gorm.io/gorm"
)

type repository struct {
	dbService *persistence.DbService
}

func NewRepository(dbService *persistence.DbService) event_repository_contract.Repository {
	return &repository{
		dbService: dbService,
	}
}

func (rp *repository) GetByMessageID(ctx context.Context, messageID string) (*entities.Event, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	event := new(entities.Event)

	result := tx.Where("message_id = ?", messageID).First(event)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, app_error.New(http.StatusNotFound, fmt.Sprintf("event with message id %s not found", messageID))
		}

		return nil, result.Error
	}

	return event, nil
}

func (rp *repository) GetByIDsAndDateTime(ctx context.Context, eventData *entities.Event) (*entities.Event, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	event := new(entities.Event)

	query := tx.Where("user_id = ?", eventData.UserID)
	query = query.Where("event_type = ?", eventData.EventType)
	query = query.Where("outcome IS NULL")

	result := query.First(event)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, app_error.New(http.StatusNotFound, fmt.Sprintf("event with user id %d and event type %s not found", eventData.UserID, eventData.EventType))
		}

		return nil, result.Error
	}

	return event, nil
}

func (rp *repository) Create(ctx context.Context, event *entities.Event) (*entities.Event, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Create(event).Error; err != nil {
		return nil, err
	}

	return event, nil
}

func (rp *repository) Update(ctx context.Context, event *entities.Event) (*entities.Event, error) {
	tx := rp.dbService.Instance.WithContext(ctx)

	if err := tx.Model(event).Where("message_id = ?", event.MessageID).Updates(event).Error; err != nil {
		return nil, err
	}

	return event, nil
}
