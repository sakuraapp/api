package controller

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/markbates/goth/gothic"
	"github.com/sakuraapp/api/resource"
	"github.com/sakuraapp/shared/model"
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
		fmt.Printf("%v", err.Error())
		render.Render(w, r, resource.ErrBadRequest)
		return
	}

	fmt.Printf("User: %+v\n", extUser)

	repos := c.app.GetRepositories()
	user, err := repos.User.GetByExternalIdWithDiscriminator(extUser.UserID)

	if err != nil {
		fmt.Printf("%v", err)
		c.SendInternalError(w, r)
		return
	}

	if user == nil {
		discrim, err := repos.Discriminator.FindFreeOne(name)

		if err != nil {
			fmt.Printf("%v", err)
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
			fmt.Printf("%v", err)
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
			fmt.Printf("%v", err)
			c.SendInternalError(w, r)
			return
		}
	} else {
		user.AccessToken = accessToken
		user.RefreshToken = refreshToken
		user.Avatar = avatar

		err := repos.User.Update(user)

		if err != nil {
			fmt.Printf("%v", err)
			c.SendInternalError(w, r)
			return
		}
	}

	payload := map[string]interface{}{
		"id": user.Id,
	}

	_, t, err := c.app.GetJWT().Encode(payload)

	if err != nil {
		fmt.Printf("%v", err)
		c.SendInternalError(w, r)
		return
	}

	res := resource.NewAuthResponse(&t)
	render.Render(w, r, res)
}