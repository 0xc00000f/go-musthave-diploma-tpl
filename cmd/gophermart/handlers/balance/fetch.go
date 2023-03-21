package balance

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/storage"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/structures"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/structures/status"
)

type Fetcher interface {
	FetchUserInfo(ctx context.Context, username string) (*storage.UserInfoData, error)
}

type FetchResp []Order

type Order struct {
	OrderNumber string             `json:"number"`
	Status      status.OrderStatus `json:"status"`
	Accrual     int64              `json:"accrual,omitempty"`
	CreatedTS   string             `json:"uploaded_at"`
}

func FetchUserInfo(fetcher Fetcher) func(*gin.Context) {
	return func(c *gin.Context) {
		todoUser := "todoUser"

		data, err := fetcher.FetchUserInfo(c, todoUser)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)

			return
		}

		c.JSON(http.StatusOK, structures.UserInfoFromStorageData(*data))
	}
}
