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

	sc := configs.Store()
	storageService, err := storage.New(sc)
	if err != nil {
		logHandler.Fatal(sc.Hosts, sc.Name, sc.AuthenticationSource, err.Error())
		return
	}
	cc := configs.Cache()
	cacheService, err := cache.New(cc)
	if err != nil {
		logHandler.Fatal(cc.Hosts, cc.Name, err.Error())
		return
	}

	serviceHandler := services.New(storageService, cacheService, logHandler)
	apiHandler := api.NewHandler(serviceHandler)

	router := webgo.NewRouter(configs.Webgo(), apiHandler.Routes())
	router.Use(middleware.AccessLog)
	router.Start()
}
