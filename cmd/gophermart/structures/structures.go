package structures

import (
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/storage"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/structures/status"
)

type Order struct {
	OrderNumber string             `json:"number"`
	Status      status.OrderStatus `json:"status"`
	Accrual     float64            `json:"accrual"`
	Withdraw    float64            `json:"withdraw"`
	CreatedTS   string             `json:"uploaded_at"`
}

type UserInfo struct {
	Balance  float64 `json:"balance"`
	Withdraw float64 `json:"withdraw"`
}

func UserInfoFromStorageData(data storage.UserInfoData) UserInfo {
	return UserInfo{
		Balance:  data.Balance,
		Withdraw: data.Withdraw,
	}
}

type Accrual struct {
	OrderNumber string               `json:"order"`
	Status      status.AccrualStatus `json:"status"`
	Accrual     float64              `json:"accrual"`
}
