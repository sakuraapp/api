package controller

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sakuraapp/api/middleware"
	"github.com/sakuraapp/api/response"
	"github.com/sakuraapp/shared/model"
	"github.com/sakuraapp/shared/resource"
	"net/http"
	"strconv"
)

type RoomController struct {
	Controller
}

func (c *RoomController) Get(w http.ResponseWriter, r *http.Request)  {
	strRoomId := chi.URLParam(r, "roomId")
	roomId, err := strconv.ParseInt(strRoomId, 10, 64)

	if err != nil {
		render.Render(w, r, response.ErrBadRequest)
		return
	}

	room, err := c.app.GetRepositories().Room.Get(model.RoomId(roomId))

	if err != nil {
		render.Render(w, r, response.ErrInternalError)
	}

	response := response.NewRoomResponse(resource.NewRoom(room))
	render.Render(w, r, response)
}

func (c *RoomController) GetLatest(w http.ResponseWriter, r *http.Request) {
	rooms, err := c.app.GetRepositories().Room.GetLatest()

	if err != nil {
		render.Render(w, r, response.ErrInternalError)
		return
	}

	response := response.NewRoomListResponse(resource.NewRoomList(rooms))
	render.Render(w, r, response)
}

func (c *RoomController) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.FromContext(r.Context())

	roomRepo := c.app.GetRepositories().Room
	room, err := roomRepo.GetByOwnerId(user.Id)

	if err != nil {
		render.Render(w, r, response.ErrInternalError)
		return
	}

	if room == nil {
		room = &model.Room{
			Name: fmt.Sprintf("%s's room", user.Username),
			OwnerId: user.Id,
			Private: false,
		}

		err = roomRepo.Create(room)

		if err != nil {
			render.Render(w, r, response.ErrInternalError)
			return
		}
	}

	response := response.NewRoomResponse(resource.NewRoom(room))
	render.Render(w, r, response)
}