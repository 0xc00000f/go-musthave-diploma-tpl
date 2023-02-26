package main

import (
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/config"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/service"
)

func main() {
	cfg, err := config.Load("config")
	if err != nil {
		panic("config file not found")
	}

	api := service.New(cfg)
	api.CreateHTTPEndpoints()
	api.Run()
}
