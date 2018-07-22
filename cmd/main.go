package main

import (
	"github.com/bnkamalesh/webgo"
	"github.com/bnkamalesh/webgo/middleware"

	"github.com/bnkamalesh/notes/api"
	"github.com/bnkamalesh/notes/configs"
	"github.com/bnkamalesh/notes/pkg/platform/logger"
	"github.com/bnkamalesh/notes/pkg/platform/storage"
	"github.com/bnkamalesh/notes/pkg/services"
)

func main() {
	logHandler := logger.New()

	storageService, err := storage.New(configs.Store())
	if err != nil {
		logHandler.Fatal(err.Error())
		return
	}

	serviceHandler := services.New(storageService, logHandler)
	apiHandler := api.NewHandler(serviceHandler)

	router := webgo.NewRouter(configs.Webgo(), apiHandler.Routes())
	router.Use(middleware.AccessLog)
	router.Start()
}
