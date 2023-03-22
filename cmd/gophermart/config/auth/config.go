package auth

import (
	"time"
)

type Config struct {
	JWTKey      string        `mapstructure:"jwt-key"`
	JWTIssuer   string        `mapstructure:"jwt-issuer"`
	JWTTokenTTL time.Duration `mapstructure:"jwt-token-ttl"`
}
