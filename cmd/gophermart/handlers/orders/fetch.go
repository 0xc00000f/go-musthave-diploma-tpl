package orders

import (
	"context"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/storage"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/structures"
)

type Fetcher interface {
	Fetch(ctx context.Context, numbers []string) (storage.OrderDataMap, error)
	FetchByUser(ctx context.Context, username string) (storage.OrderDataMap, error)
}

type Order struct {
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

func respFromOrderDataMap(orders storage.OrderDataMap) []structures.Order {
	var resp []structures.Order //nolint:prealloc

	for _, order := range orders {
		resp = append(resp, structures.OrderFromStorageData(*order))
	}

	sort.Slice(resp, func(i, j int) bool {
		return resp[i].CreatedTS < resp[j].CreatedTS
	})

	return resp
}
