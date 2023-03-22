package ctxdata

import (
	"github.com/gin-gonic/gin"
)

type ctxKey string

const (
	ctxKeyUsername ctxKey = "username"
)

func GetUsername(ctx *gin.Context) (string, bool) {
	value, exists := ctx.Get(string(ctxKeyUsername))
	if exists {
		username, ok := value.(string)
		if ok {
			return username, true
		}

		return "", false
	}

	return "", false
}

func SetUsername(ctx *gin.Context, username string) {
	ctx.Set(string(ctxKeyUsername), username)
}
