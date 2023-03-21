package structures

import (
	"time"

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/storage"
)

type Order struct {
	OrderNumber string `json:"number"`
	Status      string `json:"status"`
	Accrual     int64  `json:"accrual"`
	Withdraw    int64  `json:"withdraw"`
	CreatedTS   string `json:"uploaded_at"`
}

func OrderFromStorageData(data storage.OrderData) Order {
	return Order{
		OrderNumber: data.OrderNumber,
		Status:      data.Status,
		Accrual:     data.Accrual,
		Withdraw:    data.Withdraw,
		CreatedTS:   time.Unix(data.CreatedTS, 0).Format(time.RFC3339),
	}
}

type AccrualStatus string

const (
	AccrualStatusRegistered AccrualStatus = "REGISTERED"
	AccrualStatusInvalid    AccrualStatus = "INVALID"
	AccrualStatusProcessing AccrualStatus = "PROCESSING"
	AccrualStatusProcessed  AccrualStatus = "PROCESSED"
)

type AccrualResponse struct {
	OrderNumber string        `json:"order"`
	Status      AccrualStatus `json:"status"`
	Accrual     float64       `json:"accrual"`
}
