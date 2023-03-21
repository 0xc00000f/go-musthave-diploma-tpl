package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/crypto"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/storage"
)

type Register interface {
	Register(ctx context.Context, user storage.UserData) error
}

type RegisterReq struct { //nolint:musttag
	Username string `query:"login" validate:"required" required:"true"`
	Password string `query:"password" validate:"required" required:"true"`
}

func RegisterUser(register Register) func(*gin.Context) {
	return func(c *gin.Context) {
		var req RegisterReq
		if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)

			return
		}

		passHash, err := crypto.HashAndSalt([]byte(req.Password))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)

			return
		}

		err = register.Register(c, storage.UserData{
			Username: req.Username,
			Password: passHash,
		})

		switch {
		case errors.Is(err, storage.ErrUserAlreadyExists):
			c.AbortWithStatus(http.StatusConflict)

			return

		case err != nil:
			c.AbortWithStatus(http.StatusInternalServerError)

			return
		}

		c.Status(http.StatusOK)
	}
}
