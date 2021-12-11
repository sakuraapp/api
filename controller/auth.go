package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/markbates/goth/gothic"
	"github.com/sakuraapp/api/resource"
	"github.com/sakuraapp/shared/model"
	log "github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v4"
	"net/http"
)

type AuthController struct {
	Controller
}

func (c *AuthController) BeginAuth(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	req := gothic.GetContextWithProvider(r, provider)

	gothic.BeginAuthHandler(w, req)
}

func (c *AuthController) CompleteAuth(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")

	req := gothic.GetContextWithProvider(r, provider)
	extUser, err := gothic.CompleteUserAuth(w, req)

	name := extUser.Name
	accessToken := null.StringFrom(extUser.AccessToken)
	refreshToken := null.StringFrom(extUser.RefreshToken)
	avatar := null.StringFrom(extUser.AvatarURL)

	if err != nil {
		// todo: handle this error?
		log.WithError(err).Error("Failed to complete user auth")
		render.Render(w, r, resource.ErrBadRequest)
		return
	}

	log.Debugf("User: %+v", extUser)

	repos := c.app.GetRepositories()
	user, err := repos.User.GetByExternalIdWithDiscriminator(extUser.UserID)

	if err != nil {
		log.
			WithFields(log.Fields{
				"provider": extUser.Provider,
				"external_user_id": extUser.UserID,
			}).
			WithError(err).
			Error("Failed to get user by external id")

		c.SendInternalError(w, r)
		return
	}

	if user == nil {
		discrim, err := repos.Discriminator.FindFreeOne(name)

		if err != nil {
			log.
				WithField("name", name).
				WithError(err).
				Error("Failed to find a free discriminator")

			c.SendInternalError(w, r)
			return
		}

		if discrim == nil {
			render.Render(w, r, resource.ErrTooManyUsers)
			return
		}

		user = &model.User{
			Username: name,
			Avatar: avatar,
			Provider: extUser.Provider,
			AccessToken: accessToken,
			RefreshToken: refreshToken,
			ExternalUserID: null.StringFrom(extUser.UserID),
			Discriminator: null.StringFromPtr(discrim),
		}

		err = repos.User.Create(user)

		if err != nil {
			log.WithError(err).Error("Failed to create user")
			c.SendInternalError(w, r)
			return
		}

		discriminator := &model.Discriminator{
			Name: name,
			Value: *discrim,
			OwnerId: user.Id,
		}

		err = repos.Discriminator.Create(discriminator)

		if err != nil {
			log.
				WithField("discriminator", discrim).
				WithError(err).
				Error("Failed to insert discriminator")

			c.SendInternalError(w, r)
			return
		}
	} else {
		user.AccessToken = accessToken
		user.RefreshToken = refreshToken
		user.Avatar = avatar

		err = repos.User.Update(user)

		if err != nil {
			log.
				WithField("user_id", user.Id).
				WithError(err).
				Error("Failed to update user")

			c.SendInternalError(w, r)
			return
		}
	}

	payload := map[string]interface{}{
		"id": user.Id,
	}

	_, t, err := c.app.GetJWT().Encode(payload)

	if err != nil {
		log.WithError(err).Error("Failed to encode JWT")
		c.SendInternalError(w, r)
		return
	}

	res := resource.NewAuthResponse(&t)
	render.Render(w, r, res)
}