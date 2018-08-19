package main

import (
	"github.com/bnkamalesh/webgo"
	"github.com/bnkamalesh/webgo/middleware"

	"github.com/bnkamalesh/notes/api"
	"github.com/bnkamalesh/notes/configs"
	"github.com/bnkamalesh/notes/pkg/platform/cache"
	"github.com/bnkamalesh/notes/pkg/platform/logger"
	"github.com/bnkamalesh/notes/pkg/platform/storage"
	"github.com/bnkamalesh/notes/pkg/services"
)

func main() {
	logHandler := logger.New(configs.Logs())

	storageService, err := storage.New(configs.Store())
	if err != nil {
		logHandler.Fatal(err.Error())
		return
	}
	cacheService, err := cache.New(configs.Cache())
	if err != nil {
		logHandler.Fatal(err.Error())
		return
	}

	serviceHandler := services.New(storageService, cacheService, logHandler)
	apiHandler := api.NewHandler(serviceHandler)

	router := webgo.NewRouter(configs.Webgo(), apiHandler.Routes())
	router.Use(middleware.AccessLog)
	router.Start()
}
