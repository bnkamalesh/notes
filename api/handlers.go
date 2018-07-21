package api

import (
	"net/http"

	"github.com/bnkamalesh/webgo"
)

// Home is the home page handler
func (h *Handler) Home(rw http.ResponseWriter, req *http.Request) {
	webgo.R200(rw, "Hello world")
}
