package users

import (
	"github.com/bnkamalesh/notes/pkg/items"
	"github.com/bnkamalesh/notes/pkg/platform/logger"
	"github.com/bnkamalesh/notes/pkg/platform/storage"
)

// Service holds all the dependencies of items
type Service struct {
	store  storage.Service
	items  items.Service
	logger logger.Service
}

// NewService returns a new instance of Service with all the dependencies initialized
func NewService(ss storage.Service, l logger.Service, i items.Service) Service {
	return Service{
		store:  ss,
		logger: l,
		items:  i,
	}
}
