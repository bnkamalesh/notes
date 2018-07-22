// Package api serves all the API endpoints of the app
package api

import "github.com/bnkamalesh/notes/pkg/services"

// Handler holds all the services which will serve the endpoints
type Handler struct {
	Services services.Handler
}

// NewHandler returns a handler instance with all the services initialized
func NewHandler(s services.Handler) Handler {
	return Handler{
		Services: s,
	}
}
