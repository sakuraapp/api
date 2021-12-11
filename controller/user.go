package controller

import (
	"github.com/go-chi/render"
	"github.com/sakuraapp/api/middleware"
	apiResource "github.com/sakuraapp/api/resource"
	"github.com/sakuraapp/shared/resource"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type UserController struct {
	Controller
}

func (c *UserController) GetMyUser(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserFromContext(r.Context())
	userResource := apiResource.NewUserResponse(resource.NewUser(user))

	log.Debugf("%+v", userResource)
	render.Render(w, r, userResource)
}