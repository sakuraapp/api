package controller

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
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

const defaultQueueLimit = 20
const maxQueueLimit = 50

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

	ctx := r.Context()
	room, err := c.app.GetRepositories().Room.Get(ctx, model.RoomId(roomId))

	if err != nil {
		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	if room == nil {
		render.Render(w, r, apiResource.ErrNotFound)
		return
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
	user := middleware.UserFromContext(r.Context())

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

func (c *RoomController) Update(w http.ResponseWriter, r *http.Request) {
	strRoomId := chi.URLParam(r, "roomId")
	roomId, err := strconv.Atoi(strRoomId)

	if err != nil {
		render.Render(w, r, apiResource.ErrBadRequest)
		return
	}

	data := &apiResource.RoomUpdateRequest{}
	err = render.Bind(r, data)

	if err != nil || len(data.Name) == 0 {
		render.Render(w, r, apiResource.ErrBadRequest)
		return
	}

	ctx := r.Context()

	roomRepo := c.app.GetRepositories().Room
	room, err := roomRepo.Get(ctx, model.RoomId(roomId))

	if err != nil {
		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	if room == nil {
		render.Render(w, r, apiResource.ErrNotFound)
		return
	}

	user := middleware.UserFromContext(ctx)

	// todo: add MANAGE_ROOM permission - need to separate permissions from the session (and attach it to the user themself) and possibly add roles

	if user.Id != room.OwnerId {
		render.Render(w, r, apiResource.ErrForbidden)
		return
	}

	room.Name = data.Name
	room.Private = data.Private

	err = roomRepo.UpdateInfo(room)

	if err != nil {
		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	cacheKey := fmt.Sprintf(constant.RoomCacheFmt, roomId)
	err = c.app.GetCache().Delete(ctx, cacheKey)

	if err != nil {
		fmt.Printf("Error deleting room cache: %v", err.Error())
		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	message := resource.ServerMessage{
		Data: resource.BuildPacket(opcode.UpdateRoom, data),
	}

	bytes, err := msgpack.Marshal(message)

	if err != nil {
		fmt.Printf("Serialization Error: %v", err.Error())
		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	roomKey := fmt.Sprintf(constant.RoomFmt, roomId)

	rdb := c.app.GetRedis()
	err = rdb.Publish(ctx, roomKey, bytes).Err()

	if err != nil {
		fmt.Printf("Error publishing message: %v", err.Error())
		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	render.Render(w, r, apiResource.NewResponse(200))
}

func (c *RoomController) SendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := middleware.UserFromContext(ctx)
	sess := middleware.SessionFromContext(ctx)

	data := &apiResource.MessageRequest{}
	err := render.Bind(r, data)

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
			IgnoredSessionIds: map[string]bool{sess.Id: true},
		},
		Data: resource.BuildPacket(opcode.SendMessage, msg),
	}

	bytes, err := msgpack.Marshal(message)

	if err != nil {
		fmt.Printf("Serialization Error: %v", err.Error())
		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	rdb := c.app.GetRedis()
	err = rdb.Publish(ctx, roomKey, bytes).Err()

	if err != nil {
		fmt.Printf("Error publishing message: %v", err.Error())
		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	res := apiResource.NewMessageResponse(&msg)
	render.Render(w, r, res)
}

func (c *RoomController) GetQueue(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	limit, err := strconv.ParseInt(q.Get("limit"), 10, 64)

	if err != nil {
		limit = defaultQueueLimit
	}

	after := q.Get("after")
	start := int64(0)

	ctx := r.Context()
	rdb := c.app.GetRedis()

	sess := middleware.SessionFromContext(ctx)
	roomId := sess.RoomId

	queueKey := fmt.Sprintf(constant.RoomQueueFmt, roomId)

	if after != "" {
		start, err = rdb.LPos(ctx, queueKey, after, redis.LPosArgs{}).Result()

		if err != nil {
			start = 0
			fmt.Printf("Error reading item at index %v: %v", after, err.Error())
		} else {
			start += 1 // start 1 element after the specified one
		}
	}

	ids, err := rdb.LRange(ctx, queueKey, start, start + limit - 1).Result()

	if err != nil {
		fmt.Printf("Error fetching queue: %v", err.Error())
		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	numItems := len(ids)
	queueItemsKey := fmt.Sprintf(constant.RoomQueueItemsFmt, roomId)

	var items []resource.MediaItem

	if numItems > 0 {
		rawItems, err := rdb.HMGet(ctx, queueItemsKey, ids...).Result()

		if err != nil {
			fmt.Printf("Error fetching queue items: %v", err.Error())
			render.Render(w, r, apiResource.ErrInternalError)
			return
		}

		items = make([]resource.MediaItem, numItems)

		for i, rawItem := range rawItems {
			strItem, ok := rawItem.(string)

			if ok {
				byteItem := []byte(strItem)
				err = items[i].UnmarshalBinary(byteItem)

				if err != nil {
					fmt.Printf("Deformed queue item: %v\n", err.Error())
				}
			} else {
				fmt.Printf("Deformed queue item: %v\n", rawItem)
			}
		}
	} else {
		items = []resource.MediaItem{}
	}

	res := apiResource.NewQueueResponse(items)
	render.Render(w, r, res)
}