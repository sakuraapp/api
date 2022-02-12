package resource

import (
	"github.com/sakuraapp/shared/pkg/resource"
	"net/http"
)

type AuthResponse struct {
	resource.Response
	Token *string `json:"token"`
}

func NewAuthResponse(token *string) *AuthResponse {
	return &AuthResponse{
		Response: resource.Response{Status: 200},
		Token:    token,
	}
}

var ErrTooManyUsers = resource.NewError(http.StatusConflict, resource.NewErrorMessage("Too many users with this username"))