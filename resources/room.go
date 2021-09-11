package resources

import (
	"github.com/go-chi/render"
	"github.com/sakuraapp/api/models"
	"net/http"
)

type Room struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	Owner *User `json:"owner"`
	Private bool `json:"private"`
}

func NewRoom(room *models.Room) *Room {
	return &Room{
		room.Id,
		room.Name,
		NewUser(room.Owner),
		room.Private,
	}
}

func NewRoomList(rooms []models.Room) []*Room {
	list := make([]*Room, len(rooms))

	for i, v := range rooms {
		list[i] = NewRoom(&v)
	}

	return list
}

type RoomResponse struct {
	Response
	Room *Room `json:"room,omitempty"`
}

func (res *RoomResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, res.Status)

	return nil
}

func NewRoomResponse(room *Room) *RoomResponse {
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
	Rooms []*Room `json:"rooms"`
}

func (res *RoomListResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, res.Status)

	return nil
}

func NewRoomListResponse(rooms []*Room) *RoomListResponse {
	return &RoomListResponse{
		Response{200},
		rooms,
	}
}