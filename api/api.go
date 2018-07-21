// Package api serves all the API endpoints of the app
package api

// Handler holds all the services which will serve the endpoints
type Handler struct {
}

// NewHandler returns a handler instance with all the services initialized
func NewHandler() Handler {
	return Handler{}
}
