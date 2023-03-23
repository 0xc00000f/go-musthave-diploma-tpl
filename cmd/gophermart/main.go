package main

import (
	"context"

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/config"
	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/service"
)

func main() {
	cfg := config.MustLoad("config").ParseFlags()

	api := service.New(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go api.RunAccrual(ctx)

	api.CreateHTTPEndpoints()
	api.Run()
}
