package pg

import (
	"time"
)

type Config struct {
	Dsn  string     `mapstructure:"DATABASE_URI"`
	Pool PoolConfig `mapstructure:"pool"`
}

type PoolConfig struct {
	MaxOpen     int           `mapstructure:"max-open"`
	MaxIdle     int           `mapstructure:"max-idle"`
	MaxIdleTime time.Duration `mapstructure:"max-idle-time"`
	MaxLifetime time.Duration `mapstructure:"max-lifetime"`
}
