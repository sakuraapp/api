package responses

import (
	"github.com/go-chi/render"
	"net/http"
)

type ErrResponse struct {
	Response
	Message *string `json:"message,omitempty"`
}

func (res *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, res.Status)

	return nil
}

func NewError(status int, message *string) *ErrResponse {
	return &ErrResponse{
		Response{status},
		message,
	}
}

func NewErrorMessage(message string) *string {
	return &message
}

var ErrBadRequest = NewError(http.StatusBadRequest, nil)
var ErrInternalError = NewError(http.StatusInternalServerError, nil)
var ErrUnauthorized = NewError(http.StatusUnauthorized, nil)