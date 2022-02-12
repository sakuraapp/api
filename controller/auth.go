package controller

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/markbates/goth/gothic"
	apiResource "github.com/sakuraapp/api/resource"
	"github.com/sakuraapp/shared/pkg/model"
	resource "github.com/sakuraapp/shared/pkg/resource"
	log "github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v4"
	"io"
	"net/http"
	"path/filepath"
)

const formMemoryLimit = 32 << 20

type AuthController struct {
	Controller
}

func (c *AuthController) BeginAuth(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	req := gothic.GetContextWithProvider(r, provider)

	gothic.BeginAuthHandler(w, req)
}

func (c *AuthController) HandleCallback(w http.ResponseWriter, r *http.Request) {
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
		render.Render(w, r, apiResource.ErrBadRequest)
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

	userId := model.UserId(0)

	if user != nil {
		userId = user.Id
	}

	if user == nil || user.Discriminator.IsZero() {
		user = &model.User{
			Id: userId,
			Username: name,
			Provider: extUser.Provider,
			AccessToken: accessToken,
			RefreshToken: refreshToken,
			ExternalUserID: null.StringFrom(extUser.UserID),
			Discriminator: null.StringFromPtr(nil),
		}

		var b []byte
		b, err = json.Marshal(user)

		if err != nil {
			log.WithError(err).Error("Failed to serialize user")
			c.SendInternalError(w, r)
			return
		}

		err = gothic.StoreInSession("user", string(b), r, w)

		if err != nil {
			log.WithError(err).Error("Failed to store user in session")
			c.SendInternalError(w, r)
			return
		}

		res := apiResource.NewUserResponse(&resource.User{
			Username: name,
			Avatar: avatar,
		})

		render.Render(w, r, res)
	} else {
		user.AccessToken = accessToken
		user.RefreshToken = refreshToken

		err = repos.User.Update(user)

		if err != nil {
			log.
				WithField("user_id", user.Id).
				WithError(err).
				Error("Failed to update user")

			c.SendInternalError(w, r)
			return
		}

		c.handleAuthSuccess(user.Id, w, r)
	}
}

func (c *AuthController) CompleteAuth(w http.ResponseWriter, r *http.Request) {
	strUser, err := gothic.GetFromSession("user", r)

	if err != nil {
		log.WithError(err).Error("Failed to get user from the session")
		render.Render(w, r, apiResource.ErrForbidden)
		return
	}

	var user model.User
	err = json.Unmarshal([]byte(strUser), &user)

	if err != nil {
		log.WithError(err).Error("Failed to parse user data in the session")
		c.SendInternalError(w, r)
		return
	}

	err = r.ParseMultipartForm(formMemoryLimit)

	if err != nil {
		log.WithError(err).Error("Failed to parse form data")
		render.Render(w, r, apiResource.ErrBadRequest)
		return
	}

	username := r.FormValue("username")

	if username == "" {
		render.Render(w, r, apiResource.ErrBadRequest)
		return
	}

	avatar, avHeader, err := r.FormFile("avatar")

	if err != nil && err != http.ErrMissingFile {
		log.WithError(err).Error("Failed to read avatar file")
		render.Render(w, r, apiResource.ErrBadRequest)
		return
	}

	user.Username = username

	fmt.Printf("%+v\n", avatar)

	repos := c.app.GetRepositories()
	discrim, err := repos.Discriminator.FindFreeOne(username)

	if err != nil {
		log.
			WithField("name", username).
			WithError(err).
			Error("Failed to find a free discriminator")

		c.SendInternalError(w, r)
		return
	}

	if discrim == nil {
		render.Render(w, r, apiResource.ErrTooManyUsers)
		return
	}

	if user.Id == 0 {
		err = repos.User.Create(&user)
	} else {
		err = repos.User.Update(&user)
	}

	if err != nil {
		log.WithError(err).Error("Failed to create user")
		c.SendInternalError(w, r)
		return
	}

	discrim.OwnerId = user.Id

	if discrim.Id == 0 {
		err = repos.Discriminator.Create(discrim)
	} else {
		err = repos.Discriminator.UpdateOwnerID(discrim)
	}

	if err != nil {
		log.
			WithField("user_id", user.Id).
			WithField("discriminator", discrim).
			WithError(err).
			Error("Failed to insert discriminator")

		c.SendInternalError(w, r)
		return
	}

	if avatar != nil {
		// we have to read the first 512 bytes of the image to determine its mimetype
		buff := make([]byte, 512)
		_, err = avatar.Read(buff)

		if err != nil {
			log.WithError(err).Error("Failed to read avatar file")
			c.SendInternalError(w, r)
			return
		}

		fileType := http.DetectContentType(buff)

		if fileType != "image/jpeg" && fileType != "image/png" {
			render.Render(w, r, apiResource.ErrBadRequest)
			return
		}

		_, err = avatar.Seek(0, io.SeekStart)

		if err != nil {
			log.WithError(err).Error("Failed to seek avatar file")
			c.SendInternalError(w, r)
			return
		}

		ext := filepath.Ext(avHeader.Filename)
		err = repos.User.SetAvatar(user.Id, avatar, ext)

		if err != nil {
			log.WithError(err).Error("Failed to set user avatar")
			c.SendInternalError(w, r)
			return
		}
	}

	c.handleAuthSuccess(user.Id, w, r)
}

func (c *AuthController) handleAuthSuccess(userId model.UserId, w http.ResponseWriter, r *http.Request) {
	err := gothic.Logout(w, r)

	if err != nil {
		log.WithError(err).Error("Failed to destroy login session")
		c.SendInternalError(w, r)
		return
	}

	payload := map[string]interface{}{
		"id": userId,
	}

	_, t, err := c.app.GetJWT().Encode(payload)

	if err != nil {
		log.WithError(err).Error("Failed to encode JWT")
		c.SendInternalError(w, r)
		return
	}

	res := apiResource.NewAuthResponse(&t)
	render.Render(w, r, res)
}