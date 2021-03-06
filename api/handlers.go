package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/bnkamalesh/notes/pkg/users"
	"github.com/bnkamalesh/webgo"
)

func paginationParams(req *http.Request) (int, int) {
	start := strings.TrimSpace(req.URL.Query().Get("start"))
	limit := strings.TrimSpace(req.URL.Query().Get("start"))
	startInt := 0
	limitInt := 0
	if start != "" {
		startInt, _ = strconv.Atoi(start)
	}
	if limit != "" {
		limitInt, _ = strconv.Atoi(limit)
	}
	return startInt, limitInt
}

// Home is the home page handler
func (h *Handler) Home(rw http.ResponseWriter, req *http.Request) {
	webgo.R200(rw, map[string]string{
		"version": "0.5.0",
	})
}

func (h *Handler) userSignup(rw http.ResponseWriter, req *http.Request) {
	input := make(map[string]string, 3)
	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		webgo.R400(rw, err.Error())
		return
	}
	user, err := users.New(input)
	if err != nil {
		webgo.R400(rw, err.Error())
		return
	}

	services := h.Services
	user, err = services.Users.Create(*user)
	if err != nil {
		switch err {
		case users.ErrUsrExists:
			{
				webgo.SendResponse(rw, err.Error(), 409)
			}
		default:
			webgo.R400(rw, err.Error())
		}
		return
	}
	webgo.R200(rw, user)
}

func (h *Handler) userLogin(rw http.ResponseWriter, req *http.Request) {
	input := make(map[string]string, 2)
	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		webgo.R400(rw, err.Error())
		return
	}
	services := h.Services
	user, err := services.Users.Authenticate(input["email"], input["password"], req.RemoteAddr)
	if err != nil {
		switch err {
		case users.ErrUsrNotExists:
			{
				webgo.R404(rw, err.Error())
				return
			}
		}
		webgo.R400(rw, err.Error())
		return
	}
	webgo.R200(rw, user)
}

// userItems returns the items owned by the logged in user
func (h *Handler) userItems(rw http.ResponseWriter, req *http.Request) {
	user := getUser(req)
	if user == nil {
		webgo.R403(rw, "Unidentified user")
		return
	}

	services := h.Services
	start, limit := paginationParams(req)
	items, err := services.Users.Items(user, start, limit)
	if err != nil {
		webgo.R400(rw, err.Error())
		return
	}
	webgo.R200(rw, items)
}

// userCreateItem creates a new item for the user
func (h *Handler) userCreateItem(rw http.ResponseWriter, req *http.Request) {
	user := getUser(req)
	if user == nil {
		webgo.R403(rw, "Unidentified user")
		return
	}
	input := make(map[string]string, 0)
	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		webgo.R400(rw, err.Error())
		return
	}
	services := h.Services
	item, err := services.Users.CreateItem(user, input)
	if err != nil {
		webgo.R400(rw, err.Error())
		return
	}
	webgo.R200(rw, item)
}

// userReadItem reads an existing item for the user
func (h *Handler) userReadItem(rw http.ResponseWriter, req *http.Request) {
	user := getUser(req)
	if user == nil {
		webgo.R403(rw, "Unidentified user")
		return
	}
	wctx := webgo.Context(req)
	id := wctx.Params["id"]
	services := h.Services
	item, err := services.Users.Item(user, id)
	if err != nil {
		webgo.R400(rw, err.Error())
		return
	}
	webgo.R200(rw, item)
}

// userUpdateItem updates an existing item for the user
func (h *Handler) userUpdateItem(rw http.ResponseWriter, req *http.Request) {
	user := getUser(req)
	if user == nil {
		webgo.R403(rw, "Unidentified user")
		return
	}
	input := make(map[string]string, 0)
	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		webgo.R400(rw, err.Error())
		return
	}
	wctx := webgo.Context(req)
	id := wctx.Params["id"]
	services := h.Services
	item, err := services.Users.UpdateItem(user, id, input)
	if err != nil {
		webgo.R400(rw, err.Error())
		return
	}
	webgo.R200(rw, item)
}

// userDeleteItem delets an item for the user
func (h *Handler) userDeleteItem(rw http.ResponseWriter, req *http.Request) {
	user := getUser(req)
	if user == nil {
		webgo.R403(rw, "Unidentified user")
		return
	}
	wctx := webgo.Context(req)
	id := wctx.Params["id"]
	services := h.Services

	item, err := services.Users.DeleteItem(user, id)
	if err != nil {
		webgo.R400(rw, err.Error())
		return
	}
	webgo.R200(rw, item)
}
