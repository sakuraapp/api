package controller

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/sakuraapp/api/middleware"
	apiResource "github.com/sakuraapp/api/resource"
	"github.com/sakuraapp/shared/resource"
	"net/http"
)

type UserController struct {
	Controller
}

func (c *UserController) GetMyUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	userResource := apiResource.NewUserResponse(resource.NewUser(user))

	fmt.Printf("%+v", userResource)
	render.Render(w, r, userResource)
}