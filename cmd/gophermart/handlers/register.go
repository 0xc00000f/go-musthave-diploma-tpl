package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/storage"
)

type Register interface {
	Register(ctx context.Context, user storage.User) error
}

type FetchReq struct { //nolint:musttag
	Username string `query:"login" validate:"required" required:"true"`
	Password string `query:"password" validate:"required" required:"true"`
}

func RegisterUser(register Register) func(*gin.Context) {
	return func(c *gin.Context) {
		var req FetchReq
		if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)

			return
		}

		err := register.Register(c, storage.User{Username: req.Username, Password: req.Password})

		switch {
		case errors.Is(err, storage.ErrUserAlreadyExists):
			c.AbortWithStatus(http.StatusConflict)

			return

		case err != nil:
			c.AbortWithStatus(http.StatusInternalServerError)

			return
		}
	}
}
