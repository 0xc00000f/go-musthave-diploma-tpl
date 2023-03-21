package orders

import (
	"context"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/storage"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/structures/status"
)

type Fetcher interface {
	Fetch(ctx context.Context, numbers []string) (storage.OrderDataMap, error)
	FetchByUser(ctx context.Context, username string) (storage.OrderDataMap, error)
}

type FetchResp []Order

type Order struct {
	OrderNumber string             `json:"number"`
	Status      status.OrderStatus `json:"status"`
	Accrual     int64              `json:"accrual,omitempty"`
	CreatedTS   string             `json:"uploaded_at"`
}

func FetchOrder(fetcher Fetcher) func(*gin.Context) {
	return func(c *gin.Context) {
		todoUser := "todoUser"

		orders, err := fetcher.FetchByUser(c, todoUser)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)

			return
		}

		if len(orders) == 0 {
			c.AbortWithStatus(http.StatusNoContent)

			return
		}

		c.JSON(http.StatusOK, respFromOrderDataMap(orders))
	}
}

func respFromOrderDataMap(data storage.OrderDataMap) FetchResp {
	var resp FetchResp //nolint:prealloc

	for _, order := range data {
		resp = append(resp, Order{
			OrderNumber: order.OrderNumber,
			Status:      order.Status,
			Accrual:     order.Accrual,
			CreatedTS:   time.Unix(order.CreatedTS, 0).Format(time.RFC3339),
		})
	}

	sort.Slice(resp, func(i, j int) bool {
		return resp[i].CreatedTS < resp[j].CreatedTS
	})

	return resp
}
