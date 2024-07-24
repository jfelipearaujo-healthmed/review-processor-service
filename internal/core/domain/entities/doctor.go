package entities

import (
	"gorm.io/gorm"
)

type Doctor struct {
	gorm.Model

	MedicalID     string  `json:"medical_id" gorm:"uniqueIndex"`
	Specialty     string  `json:"specialty"`
	Price         float64 `json:"price"`
	AvgRating     float64 `json:"avg_rating"`
	TotalPatients int     `json:"total_patients"`

	UserID uint `json:"user_id"`
}
