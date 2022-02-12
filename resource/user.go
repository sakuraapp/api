package resource

import (
	"github.com/go-chi/render"
	"github.com/sakuraapp/shared/pkg/resource"
	"net/http"
)

type UserResponse struct {
	Response
	User *resource.User `json:"user,omitempty"`
}

func (res *UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, res.Status)

	return nil
}

func NewUserResponse(user *resource.User) *UserResponse {
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