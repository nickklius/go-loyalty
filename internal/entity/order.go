package entity

import (
	"time"
)

type Order struct {
	ID         string    `json:"id,omitempty"`
	UserID     string    `json:"user_id,omitempty"`
	Number     string    `json:"number,omitempty"`
	Status     string    `json:"status,omitempty"`
	UploadedAt time.Time `json:"uploaded_at,omitempty"`
	Accrual    float64   `json:"accrual,omitempty"`
}
