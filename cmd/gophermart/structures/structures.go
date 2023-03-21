package structures

import (
	"time"

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/storage"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/structures/status"
)

type Order struct {
	OrderNumber string             `json:"number"`
	Status      status.OrderStatus `json:"status"`
	Accrual     int64              `json:"accrual"`
	Withdraw    int64              `json:"withdraw"`
	CreatedTS   string             `json:"uploaded_at"`
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

type UserInfo struct {
	Balance  int64 `json:"balance"`
	Withdraw int64 `json:"withdraw"`
}

func UserInfoFromStorageData(data storage.UserInfoData) UserInfo {
	return UserInfo{
		Balance:  data.Balance,
		Withdraw: data.Withdraw,
	}
}

type AccrualResponse struct {
	OrderNumber string               `json:"order"`
	Status      status.AccrualStatus `json:"status"`
	Accrual     float64              `json:"accrual"`
}
