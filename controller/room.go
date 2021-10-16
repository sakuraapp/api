package controller

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/sakuraapp/api/internal"
	"github.com/sakuraapp/api/middleware"
	apiResource "github.com/sakuraapp/api/resource"
	"github.com/sakuraapp/shared/constant"
	"github.com/sakuraapp/shared/model"
	"github.com/sakuraapp/shared/resource"
	"github.com/sakuraapp/shared/resource/opcode"
	"github.com/vmihailenco/msgpack/v5"
	"net/http"
	"strconv"
)

type RoomController struct {
	Controller
}

func (c *RoomController) Get(w http.ResponseWriter, r *http.Request)  {
	strRoomId := chi.URLParam(r, "roomId")
	roomId, err := strconv.Atoi(strRoomId)

	if err != nil {
		render.Render(w, r, apiResource.ErrBadRequest)
		return
	}

	room, err := c.app.GetRepositories().Room.Get(model.RoomId(roomId))

	if err != nil {
		render.Render(w, r, apiResource.ErrInternalError)
	}

	res := apiResource.NewRoomResponse(resource.NewRoom(room))
	render.Render(w, r, res)
}

func (c *RoomController) GetLatest(w http.ResponseWriter, r *http.Request) {
	rooms, err := c.app.GetRepositories().Room.GetLatest()

	if err != nil {
		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	res := apiResource.NewRoomListResponse(resource.NewRoomList(rooms))
	render.Render(w, r, res)
}

func (c *RoomController) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.FromContext(r.Context())

	roomRepo := c.app.GetRepositories().Room
	room, err := roomRepo.GetByOwnerId(user.Id)

	if err != nil {
		render.Render(w, r, apiResource.ErrInternalError)
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
			render.Render(w, r, apiResource.ErrInternalError)
			return
		}
	}

	res := apiResource.NewRoomResponse(resource.NewRoom(room))
	render.Render(w, r, res)
}

func (c *RoomController) SendMessage(w http.ResponseWriter, r *http.Request) {
	strRoomId := chi.URLParam(r, "roomId")
	intRoomId, err := strconv.Atoi(strRoomId)

	if err != nil {
		render.Render(w, r, apiResource.ErrBadRequest)
		return
	}

	roomId := model.RoomId(intRoomId)

	ctx := r.Context()

	user := middleware.FromContext(ctx)
	sessionId := r.Header.Get("X-Session-Id")

	if sessionId == "" {
		render.Render(w, r, apiResource.ErrForbidden)
		return
	}

	sessionKey := fmt.Sprintf(constant.SessionFmt, sessionId)

	var sess internal.Session

	rdb := c.app.GetRedis()
	err = rdb.HMGet(ctx, sessionKey, "user_id", "room_id", "node_id").Scan(&sess)

	if err != nil {
		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	if sess.UserId != user.Id || sess.RoomId != roomId {
		render.Render(w, r, apiResource.ErrForbidden)
		return
	}

	data := &apiResource.MessageRequest{}
	err = render.Bind(r, data)

	if err != nil || len(data.Content) == 0 {
		render.Render(w, r, apiResource.ErrBadRequest)
		return
	}

	roomKey := fmt.Sprintf(constant.RoomFmt, sess.RoomId)
	id := uuid.NewString()
	msg := resource.Message{
		Id: id,
		Author: user.Id,
		Content: data.Content,
	}

	message := resource.ServerMessage{
		Target: resource.MessageTarget{
			IgnoredSessionIds: map[string]bool{sessionId: true},
		},
		Data: resource.BuildPacket(opcode.SEND_MESSAGE, msg),
	}

	bytes, err := msgpack.Marshal(message)

	if err != nil {
		fmt.Printf("Serialization Error: %v", err.Error())
		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	err = rdb.Publish(ctx, roomKey, bytes).Err()

	if err != nil {
		fmt.Printf("Error publishing message: %v", err.Error())
		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	res := apiResource.NewMessageResponse(&msg)
	render.Render(w, r, res)
}