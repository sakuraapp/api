package resources

import (
	"github.com/go-chi/render"
	"github.com/sakuraapp/api/models"
	"gopkg.in/guregu/null.v4"
	"net/http"
)

type User struct {
	Id int64 `json:"id"`
	Username string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar null.String `json:"avatar"`
}

func NewUser(user *models.User) *User {
	return &User{
		Id: user.Id,
		Username: user.Username,
		Discriminator: user.Discriminator.ValueOrZero(),
		Avatar: user.Avatar,
	}
}

type UserResponse struct {
	Response
	User *User `json:"user,omitempty"`
}

func (res *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, res.Status)

	return nil
}

func NewUserResponse(user *User) *UserResponse {
	var status int

	if user != nil {
		status = 200
	} else {
		status = 404
	}

	return &UserResponse{
		Response{status},
		user,
	}
}