package resource

import (
	"github.com/go-chi/render"
	"github.com/sakuraapp/shared/pkg/resource"
	"net/http"
)

type RoomUpdateRequest struct {
	Name string `json:"name" msgpack:"name"`
	Private bool `json:"private" msgpack:"private"`
}

func (req *RoomUpdateRequest) Bind(r *http.Request) error {
	return nil
}

type RoomResponse struct {
	Response
	Room *resource.Room `json:"room,omitempty"`
}

func (res *RoomResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, res.Status)

	return nil
}

func NewRoomResponse(room *resource.Room) *RoomResponse {
	return &RoomResponse{
		Response{200},
		room,
	}
}

type RoomListResponse struct {
	Response
	Rooms []*resource.Room `json:"rooms"`
}

func (res *RoomListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, res.Status)

	return nil
}

func NewRoomListResponse(rooms []*resource.Room) *RoomListResponse {
	return &RoomListResponse{
		Response{200},
		rooms,
	}
}