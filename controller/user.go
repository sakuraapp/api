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
	ctx := r.Context()

	userId := middleware.UserIdFromContext(ctx)
	user, err := c.app.GetRepositories().User.GetWithDiscriminator(ctx, userId)

	if err != nil {
		log.
			WithField("user_id", userId).
			WithError(err).
			Error("Failed to fetch user")

		c.SendInternalError(w, r)
		return
	}

	userResource := apiResource.NewUserResponse(resource.NewUser(user))

	log.Debugf("%+v", userResource)
	render.Render(w, r, userResource)
}