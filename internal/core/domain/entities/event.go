package entities

import (
	"time"

	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/external/queue"
	"gorm.io/gorm"
)

type Event struct {
	ID uint `json:"id,omitempty" gorm:"primaryKey"`

	UserID    uint            `json:"user_id,omitempty"`
	MessageID string          `json:"message_id,omitempty"`
	EventType queue.EventType `json:"event_type,omitempty"`
	Data      string          `json:"data,omitempty"`
	Outcome   *string         `json:"outcome,omitempty"`

	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	Appointment *Appointment `json:"appointment,omitempty" gorm:"foreignKey:EventID"`
}
