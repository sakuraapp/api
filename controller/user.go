package controller

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/sakuraapp/api/middleware"
	"github.com/sakuraapp/api/response"
	"github.com/sakuraapp/shared/resource"
	"net/http"
)

type UserController struct {
	Controller
}

func (c *UserController) GetMyUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.FromContext(r.Context())
	userResource := response.NewUserResponse(resource.NewUser(user))

	fmt.Printf("%+v", userResource)
	render.Render(w, r, userResource)
}