package withdraw

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/auth"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/handlers/orders"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/storage"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/structures/status"
	"github.com/0xc00000f/go-musthave-diploma-tpl/lib/luhn"
)

type Withdrawer interface {
	orders.Creator
	orders.Fetcher
}

type RequestReq struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func Request(withdrawer Withdrawer) func(*gin.Context) {
	return func(c *gin.Context) {
		user, err := auth.GetUsername(c)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)

			return
		}

		var req RequestReq

		if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)

			return
		}

		if num, err := strconv.Atoi(req.Order); err != nil || !luhn.Valid(num) {
			c.AbortWithStatus(http.StatusUnprocessableEntity)

			return
		}

		orders, err := withdrawer.Fetch(c, []string{req.Order})
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)

			return
		}

		if _, ok := orders[req.Order]; ok {
			c.AbortWithStatus(http.StatusConflict)

			return
		}

		_, err = withdrawer.Create(c, storage.OrderCreateData{
			Username:    user,
			OrderNumber: req.Order,
			Withdraw:    req.Sum,
			Status:      status.OrderStatusNew,
		})
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)

			return
		}

		c.Status(http.StatusOK)
	}
}
