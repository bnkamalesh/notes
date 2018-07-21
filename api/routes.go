package api

import (
	"net/http"

	"github.com/bnkamalesh/webgo"
)

// Routes returns all the HTTP routes of the app
func (handler *Handler) Routes() []*webgo.Route {
	return []*webgo.Route{
		&webgo.Route{
			Name:     "home",
			Method:   http.MethodGet,
			Pattern:  "/",
			Handlers: []http.HandlerFunc{handler.Home},
		},
	}
}
