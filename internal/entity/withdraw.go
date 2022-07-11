package entity

import "time"

type Withdraw struct {
	ID          string    `json:"id,omitempty"`
	UserID      string    `json:"user_id,omitempty"`
	OrderID     string    `json:"order,omitempty"`
	Sum         float64   `json:"sum,omitempty"`
	Status      string    `json:"status,omitempty"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
}
