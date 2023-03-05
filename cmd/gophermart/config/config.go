package config

import (
	"strings"

	"github.com/spf13/viper"

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/config/pg"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/config/webserver"
)

type Config struct {
	Viper *viper.Viper

	Webserver webserver.Config `mapstructure:"webserver"`
	Pgsql     pg.Config        `mapstructure:"pgsql"`
}

func MustLoad(path string) *Config {
	var cfg Config

	v := viper.New()

	v.AddConfigPath(".")
	v.SetConfigType("yaml")
	v.SetConfigName(path)

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		panic("config file not found")
	}

	err = v.Unmarshal(&cfg)
	if err != nil {
		panic("config file unmarshalling failed")
	}

	cfg.Viper = v

	return &cfg
}
