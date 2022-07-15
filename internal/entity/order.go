package entity

import (
	"time"
)

const (
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusRegistered OrderStatus = "REGISTERED"
)

type OrderStatus string

type Order struct {
	ID          string      `json:"id,omitempty"`
	UserID      string      `json:"user_id,omitempty"`
	Number      string      `json:"number,omitempty"`
	OrderStatus OrderStatus `json:"status,omitempty"`
	UploadedAt  time.Time   `json:"uploaded_at,omitempty"`
	Accrual     float64     `json:"accrual,omitempty"`
}
