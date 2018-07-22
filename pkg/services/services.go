// Package services defines all the services
package services

import (
	"github.com/bnkamalesh/notes/pkg/items"
	"github.com/bnkamalesh/notes/pkg/platform/logger"
	"github.com/bnkamalesh/notes/pkg/platform/storage"
)

// Handler holds all the services of the app
type Handler struct {
	Items items.Service
}

// New returns
func New(ss storage.Service, l logger.Service) Handler {
	iS := items.NewService(ss, l)
	return Handler{
		Items: iS,
	}
}
