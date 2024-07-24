package entities

import (
	"time"

	"gorm.io/gorm"
)

type Status string

func (s Status) String() string {
	return string(s)
}

const (
	ScheduleInAnalysis     Status = "schedule_in_analysis"
	ReScheduleInAnalysis   Status = "re_schedule_in_analysis"
	WaitingForConfirmation Status = "waiting_for_confirmation"
	Confirmed              Status = "confirmed"
	InProgress             Status = "in_progress"
	Concluded              Status = "concluded"
	Cancelled              Status = "cancelled"
)

type Appointment struct {
	ID uint `json:"id,omitempty" gorm:"primaryKey"`

	ScheduleID      uint       `json:"schedule_id,omitempty"`
	PatientID       uint       `json:"patient_id,omitempty"`
	DoctorID        uint       `json:"doctor_id,omitempty"`
	DateTime        time.Time  `json:"date_time,omitempty"`
	Status          Status     `json:"status,omitempty"`
	StartedAt       *time.Time `json:"started_at,omitempty"`
	EndedAt         *time.Time `json:"ended_at,omitempty"`
	ConfirmedAt     *time.Time `json:"confirmed_at,omitempty"`
	CancelledBy     *uint      `json:"cancelled_by,omitempty"`
	CancelledAt     *time.Time `json:"cancelled_at,omitempty"`
	CancelledReason *string    `json:"cancelled_reason,omitempty"`

	CreatedAt time.Time      `json:"created_at,omitempty"`
	UpdatedAt time.Time      `json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	EventID         uint `json:"event_id,omitempty"`
	FeedbackID      uint `json:"feedback_id,omitempty"`
	MedicalReportID uint `json:"medical_report_id,omitempty"`
}

func (a *Appointment) Cancel(cancelledBy uint, reason string) {
	now := time.Now()

	a.CancelledBy = &cancelledBy
	a.CancelledAt = &now
	a.CancelledReason = &reason
}
