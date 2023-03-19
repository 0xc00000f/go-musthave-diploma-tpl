package orders

import (
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/storage"
	"github.com/0xc00000f/go-musthave-diploma-tpl/lib/luhn"
)

type Creator interface {
	Create(ctx context.Context, data storage.OrderCreateData) (*storage.OrderData, error)
}

type CreateFetcher interface {
	Creator
	Fetcher
}

func CreateOrder(cf CreateFetcher) func(*gin.Context) {
	return func(c *gin.Context) {
		todoUser := "todoUser"

		if c.Request.Body == nil || c.Request.Header.Get("Content-Type") != "text/plain" {
			c.AbortWithStatus(http.StatusBadRequest)

			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)

			return
		}

		number := string(body)

		if num, err := strconv.Atoi(number); err != nil || !luhn.Valid(num) {
			c.AbortWithStatus(http.StatusUnprocessableEntity)

			return
		}

		orders, err := cf.Fetch(c, []string{number})
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)

			return
		}

		if order, ok := orders[number]; ok {
			if order.Username == todoUser {
				c.AbortWithStatus(http.StatusOK)

				return
			}

			c.AbortWithStatus(http.StatusConflict)

			return
		}

		_, err = cf.Create(c, storage.OrderCreateData{
			Username:    todoUser,
			OrderNumber: number,
		})
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)

			return
		}

		c.Status(http.StatusOK)
	}
}
