package model

import "time"

type Status string

const (
	New Status = "NEW"
	// Registrered — заказ зарегистрирован, но вознаграждение не рассчитано.
	Registrered Status = "REGISTERED"
	// Invalid — заказ не принят к расчёту, и вознаграждение не будет начислено.
	Invalid Status = "INVALID"
	// Processing — расчёт начисления в процессе.
	Processing Status = "PROCESSING"
	// Processed - расчёт начисления окончен.
	Processed Status = "PROCESSED"
)

type Order struct {
	Number     string    `json:"number"`
	Status     Status    `json:"status"`
	Accrual    float64   `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type Accrual struct {
	Order   string `json:"order"`
	Status  Status `json:"status"`
	Accrual *int   `json:"accrual,omitempty"`
}
