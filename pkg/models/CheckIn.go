package models

import (
	"time"

	"gorm.io/gorm"
)

type CheckIn struct {
	gorm.Model
	AtletaID              string    `json:"atleta_id" gorm:"type:uuid"`
	Data                  time.Time `json:"data"`
	AllenamentoCompletato bool      `json:"allenamento_completato"`
	DietaCompletata       bool      `json:"dieta_completata"`
	CardioCompletato      bool      `json:"cardio_completato"`
	LitriAcqua            float64   `json:"litri_acqua"`
	XPGained              int       `json:"xp_gained"`
}
