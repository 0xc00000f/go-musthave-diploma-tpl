package withdraw

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/storage"
)

type Fetcher interface {
	FetchByUser(ctx context.Context, username string) (storage.OrderDataMap, error)
}

type FetchResp []Order

type Order struct {
	OrderNumber string  `json:"order"`
	Withdraw    float64 `json:"sum"`
	CreatedTS   string  `json:"processed_at"`
}

func FetchUserInfo(fetcher Fetcher) func(*gin.Context) {
	return func(c *gin.Context) {
		todoUser := "todoUser"

		data, err := fetcher.FetchByUser(c, todoUser)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)

			return
		}

		resp, ok := respFromUserInfo(data)
		if !ok {
			c.AbortWithStatus(http.StatusNoContent)

			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func respFromUserInfo(data storage.OrderDataMap) (FetchResp, bool) {
	var resp FetchResp //nolint:prealloc

	for _, order := range data {
		if order.Withdraw == 0 {
			continue
		}

		resp = append(resp, Order{
			OrderNumber: order.OrderNumber,
			Withdraw:    order.Withdraw,
			CreatedTS:   time.Unix(order.CreatedTS, 0).Format(time.RFC3339),
		})
	}

	return resp, len(resp) > 0
}
