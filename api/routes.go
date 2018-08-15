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
		&webgo.Route{
			Name:     "userSignup",
			Method:   http.MethodPost,
			Pattern:  "/signup",
			Handlers: []http.HandlerFunc{handler.userSignup},
		},
		&webgo.Route{
			Name:     "userLogin",
			Method:   http.MethodPost,
			Pattern:  "/login",
			Handlers: []http.HandlerFunc{handler.userLogin},
		},
		&webgo.Route{
			Name:     "userItems",
			Method:   http.MethodGet,
			Pattern:  "/items",
			Handlers: []http.HandlerFunc{handler.mwareAuthenticate, handler.userItems},
		},
	}
}
