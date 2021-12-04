package resource

import (
	"github.com/go-chi/render"
	"github.com/sakuraapp/shared/resource"
	"net/http"
)

type QueueResponse struct {
	Response
	Items []resource.MediaItem `json:"items,omitempty"`
}

func (res *QueueResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, res.Status)

	return nil
}

func NewQueueResponse(items []resource.MediaItem) *QueueResponse {
	return &QueueResponse{
		Response{Status: 200},
		items,
	}
}
