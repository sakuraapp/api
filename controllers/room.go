package controllers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sakuraapp/api/middlewares"
	"github.com/sakuraapp/api/responses"
	"github.com/sakuraapp/shared/models"
	"github.com/sakuraapp/shared/resources"
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
		render.Render(w, r, responses.ErrBadRequest)
		return
	}

	room, err := c.app.GetRepositories().Room.Get(roomId)

	if err != nil {
		render.Render(w, r, responses.ErrInternalError)
	}

	response := responses.NewRoomResponse(resources.NewRoom(room))
	render.Render(w, r, response)
}

func (c *RoomController) GetLatest(w http.ResponseWriter, r *http.Request) {
	rooms, err := c.app.GetRepositories().Room.GetLatest()

	if err != nil {
		render.Render(w, r, responses.ErrInternalError)
		return
	}

	response := responses.NewRoomListResponse(resources.NewRoomList(rooms))
	render.Render(w, r, response)
}

func (c *RoomController) Create(w http.ResponseWriter, r *http.Request) {
	user := middlewares.FromContext(r.Context())

	roomRepo := c.app.GetRepositories().Room
	room, err := roomRepo.GetByOwnerId(user.Id)

	if err != nil {
		render.Render(w, r, responses.ErrInternalError)
		return
	}

	if room == nil {
		room = &models.Room{
			Name: fmt.Sprintf("%s's room", user.Username),
			OwnerId: user.Id,
			Private: false,
		}

		err = roomRepo.Create(room)

		if err != nil {
			render.Render(w, r, responses.ErrInternalError)
			return
		}
	}

	response := responses.NewRoomResponse(resources.NewRoom(room))
	render.Render(w, r, response)
}