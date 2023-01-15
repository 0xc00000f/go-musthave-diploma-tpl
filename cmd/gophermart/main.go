package main

import (
	"log"

	"github.com/0xc00000f/go-musthave-diploma-tpl/cmd/gophermart/handlers"
)

func main() {
	router := handlers.NewRouter()

	log.Fatal("http server down", router.Run())
}
