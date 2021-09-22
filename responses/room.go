package responses

import (
	"github.com/go-chi/render"
	"github.com/sakuraapp/shared/resources"
	"net/http"
)

type RoomResponse struct {
	Response
	Room *resources.Room `json:"room,omitempty"`
}

func (res *RoomResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, res.Status)

	return nil
}

func NewRoomResponse(room *resources.Room) *RoomResponse {
	var status int

	if room != nil {
		status = 200
	} else {
		status = 404
	}

	return &RoomResponse{
		Response{status},
		room,
	}
}

type RoomListResponse struct {
	Response
	Rooms []*resources.Room `json:"rooms"`
}

func (res *RoomListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, res.Status)

	return nil
}

func NewRoomListResponse(rooms []*resources.Room) *RoomListResponse {
	return &RoomListResponse{
		Response{200},
		rooms,
	}
}