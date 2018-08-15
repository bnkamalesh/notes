package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/bnkamalesh/notes/pkg/users"
	"github.com/bnkamalesh/webgo"
)

type userKey string

const (
	userCtxKey = userKey("user")
)

func getUser(req *http.Request) *users.User {
	u, _ := req.Context().Value(userCtxKey).(*users.User)
	return u
}

func (h *Handler) mwareAuthenticate(rw http.ResponseWriter, req *http.Request) {
	authToken := strings.TrimSpace(req.Header.Get("Authorization"))
	services := h.Services
	user, err := services.Users.AuthUser(authToken, req.RemoteAddr)
	if err != nil || authToken == "" {
		webgo.R403(rw, "Sorry, you're not authorized to access this API")
		return
	}

	reqwc := req.WithContext(
		context.WithValue(
			req.Context(),
			userCtxKey,
			user,
		),
	)
	*req = *reqwc
}
