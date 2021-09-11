package controllers

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/sakuraapp/api/middlewares"
	"github.com/sakuraapp/api/resources"
	"net/http"
)

type UserController struct {
	Controller
}

func (c *UserController) GetMyUser(w http.ResponseWriter, r *http.Request) {
	user := middlewares.FromContext(r.Context())
	userResource := resources.NewUserResponse(resources.NewUser(user))

	fmt.Printf("%+v", userResource)
	render.Render(w, r, userResource)
}