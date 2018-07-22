package items

import (
	"github.com/bnkamalesh/notes/pkg/platform/logger"
	"github.com/bnkamalesh/notes/pkg/platform/storage"
)

// Service holds all the dependencies of items
type Service struct {
	store  storage.Service
	logger logger.Service
}

// NewService returns a new instance of Service with all the dependencies initialized
func NewService(ss storage.Service, l logger.Service) Service {
	return Service{
		store: ss,
	}
}
