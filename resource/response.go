package resource

import (
	"github.com/go-chi/render"
	"net/http"
)

type Response struct {
	Status int `json:"status"`
}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, res.Status)

	return nil
}