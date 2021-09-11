package controllers

import (
	"github.com/go-chi/render"
	"github.com/sakuraapp/api/internal"
	"github.com/sakuraapp/api/resources"
	"net/http"
)

type Controller struct {
	app internal.App
}

func (c *Controller) SendInternalError(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, resources.ErrInternalError)
}

type Controllers struct {
	Auth AuthController
	User UserController
	Room RoomController
}

func Init(a internal.App) Controllers {
	return Controllers{
		Auth: AuthController{Controller{a}},
		User: UserController{Controller{a}},
		Room: RoomController{Controller{a}},
	}
}