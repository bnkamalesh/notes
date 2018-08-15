// Package services defines all the services
package services

import (
	"github.com/bnkamalesh/notes/pkg/items"
	"github.com/bnkamalesh/notes/pkg/platform/cache"
	"github.com/bnkamalesh/notes/pkg/platform/logger"
	"github.com/bnkamalesh/notes/pkg/platform/storage"
	"github.com/bnkamalesh/notes/pkg/users"
)

// Handler holds all the services of the app
type Handler struct {
	Items items.Service
	Users users.Service
}

// New returns a new Service instance with all the internal services initialized
func New(ss storage.Service, cs cache.Service, l logger.Service) Handler {
	iS := items.NewService(ss, l)
	uS := users.NewService(ss, cs, l, iS)

	return Handler{
		Items: iS,
		Users: uS,
	}
}
