package resource

import (
	"github.com/go-chi/render"
	"github.com/sakuraapp/shared/resource"
	"net/http"
)

type MessageRequest struct {
	Content string `json:"content"`
}

func (m *MessageRequest) Bind(r *http.Request) error {
	return nil
}

type MessageResponse struct {
	Response
	Message *resource.Message `json:"message,omitempty"`
}

func (res *MessageResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, res.Status)

	return nil
}

func NewMessageResponse(message *resource.Message) *MessageResponse {
	return &MessageResponse{
		Response{200},
		message,
	}
}