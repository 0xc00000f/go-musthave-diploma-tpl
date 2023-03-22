package webserver

import (
	"time"
)

type Config struct {
	Address string `mapstructure:"RUN_ADDRESS"`

	ReadHeaderTimeout time.Duration `mapstructure:"read-header-timeout"`
	ReadTimeout       time.Duration `mapstructure:"read-timeout"`
	WriteTimeout      time.Duration `mapstructure:"write-timeout"`
	IdleTimeout       time.Duration `mapstructure:"idle-timeout"`
}
