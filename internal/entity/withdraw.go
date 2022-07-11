package entity

import "time"

type Withdraw struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	OrderID     string    `json:"order"`
	Sum         float64   `json:"sum"`
	Status      string    `json:"status"`
	ProcessedAt time.Time `json:"processed_at"`
}
