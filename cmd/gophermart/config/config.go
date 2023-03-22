package config

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/config/accrual"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/config/auth"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/config/pg"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/config/webserver"
)

type Config struct {
	Viper *viper.Viper

	Webserver webserver.Config `mapstructure:"webserver"`
	Pgsql     pg.Config        `mapstructure:"pgsql"`
	Auth      auth.Config      `mapstructure:"auth"`
	Accrual   accrual.Config   `mapstructure:"accrual"`
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

func (c *Config) ParseFlags() *Config {
	pflag.String("a", "", "RUN_ADDRESS")
	pflag.String("d", "", "DATABASE_URI")
	pflag.String("r", "", "ACCRUAL_SYSTEM_ADDRESS")

	// парсим флаги
	pflag.Parse()

	// обновляем конфигурацию на основе флагов
	if addr := pflag.Lookup("a").Value.String(); addr != "" {
		c.Webserver.Address = addr
	}
	if dbURI := pflag.Lookup("d").Value.String(); dbURI != "" {
		c.Pgsql.Dsn = dbURI
	}
	if accrualAddr := pflag.Lookup("r").Value.String(); accrualAddr != "" {
		c.Accrual.Address = accrualAddr
	}

	// сохраняем конфигурацию с помощью viper
	if err := c.Viper.WriteConfig(); err != nil {
		panic(err)
	}

	return c
}
