package accrual

import (
	"time"
)

type Config struct {
	Address        string        `mapstructure:"ACCRUAL_SYSTEM_ADDRESS"`
	UpdateInterval time.Duration `mapstructure:"update-interval"`
}
