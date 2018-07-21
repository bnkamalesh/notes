package main

import (
	"github.com/bnkamalesh/webgo"
	"github.com/bnkamalesh/webgo/middleware"

	"github.com/bnkamalesh/gotodo/api"
	"github.com/bnkamalesh/gotodo/configs"
)

func main() {
	apiHandler := api.NewHandler()

	router := webgo.NewRouter(configs.Webgo(), apiHandler.Routes())
	router.Use(middleware.AccessLog)
	router.Start()
}
