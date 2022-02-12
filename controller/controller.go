package controller

import (
	"github.com/go-chi/render"
	"github.com/sakuraapp/api/internal"
	"github.com/sakuraapp/shared/pkg/resource"
	"net/http"
)

type Controller struct {
	app internal.App
}

func (c *Controller) SendInternalError(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, resource.ErrInternalError)
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