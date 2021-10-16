package response

import (
	"net/http"
)

type AuthResponse struct {
	Response
	Token *string `json:"token"`
}

func NewAuthResponse(token *string) *AuthResponse {
	return &AuthResponse{
		Response{200},
		token,
	}
}

var ErrTooManyUsers = NewError(http.StatusConflict, NewErrorMessage("Too many users with this username"))