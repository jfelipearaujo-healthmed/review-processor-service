package entities

import (
	"time"

	"gorm.io/gorm"
)

type Feedback struct {
	ID uint `json:"id,omitempty" gorm:"primaryKey"`

	AppointmentID uint    `json:"appointment_id,omitempty"`
	Rating        float64 `json:"rating,omitempty"`
	Comment       string  `json:"comment,omitempty"`

	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	Appointment *Appointment `json:"appointment,omitempty" gorm:"foreignKey:FeedbackID"`
}
