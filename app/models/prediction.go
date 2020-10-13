package models

import (
	"time"
)

// Prediction model
type Prediction struct {
	MatchID      int `gorm:"unique"`
	ModelName    string
	Predict      int
	PredictName  string
	PredictProba float64
	InsertedDate *time.Time
}
