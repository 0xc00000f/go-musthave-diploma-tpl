package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/auth/ctxdata"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/config/auth"
)

type Auth struct {
	jwtKey      []byte
	jwtIssuer   string
	jwtTokenTTL time.Duration
}

func New(config auth.Config) *Auth {
	return &Auth{
		jwtKey:      []byte(config.JWTKey),
		jwtIssuer:   config.JWTIssuer,
		jwtTokenTTL: config.JWTTokenTTL,
	}
}

type Claims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
}

func (a *Auth) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("Authorization")
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := a.parseToken(cookie)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctxdata.SetUsername(ctx, claims.Username)
		ctx.Next()
	}
}

func (a *Auth) CreateJWT(claims Claims) (string, error) {
	claims.RegisteredClaims = jwt.RegisteredClaims{ //nolint:exhaustruct
		Issuer:    a.jwtIssuer,
		ID:        uuid.New().String(),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.jwtTokenTTL)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := token.SignedString(a.jwtKey)
	if err != nil {
		return "", err //nolint:wrapcheck
	}

	return signedString, nil
}

func (a *Auth) GetTokenTTL() time.Duration {
	return a.jwtTokenTTL
}

func (a *Auth) parseToken(tokenString string) (claims *Claims, err error) {
	token, err := jwt.ParseWithClaims( //nolint:exhaustruct
		tokenString,
		&Claims{}, //nolint:exhaustruct
		func(token *jwt.Token) (any, error) {
			return a.jwtKey, nil
		},
	)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, err //nolint:wrapcheck
	}

	return claims, nil
}

func GetUsername(ctx *gin.Context) (string, error) {
	username, ok := ctxdata.GetUsername(ctx)
	if !ok {
		return "", fmt.Errorf("username not found in context")
	}

	return username, nil
}
