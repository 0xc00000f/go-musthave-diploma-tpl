package user

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/crypto"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/storage"
)

type Fetcher interface {
	Fetch(ctx context.Context, usernames []string) (storage.UserDataMap, error)
}

type AuthReq struct { //nolint:musttag
	Username string `query:"login" validate:"required" required:"true"`
	Password string `query:"password" validate:"required" required:"true"`
}

func AuthUser(fetcher Fetcher) func(*gin.Context) {
	return func(c *gin.Context) {
		var req AuthReq
		if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)

			return
		}

		users, err := fetcher.Fetch(c, []string{req.Username})
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)

			return
		}

		user, ok := users[req.Username]
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, err)

			return
		}

		isAuth := crypto.ComparePasswords(user.Password, []byte(req.Password))
		if !isAuth {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)

			return
		}

		c.Status(http.StatusOK)
	}
}
