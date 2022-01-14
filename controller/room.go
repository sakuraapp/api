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
	"github.com/sakuraapp/shared/resource/permission"
	"github.com/sakuraapp/shared/resource/role"
	log "github.com/sirupsen/logrus"
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
		log.
			WithField("room_id", roomId).
			WithError(err).
			Error("Failed to get room")

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
		log.WithError(err).Error("Failed to get latest rooms")
		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	res := apiResource.NewRoomListResponse(resource.NewRoomList(rooms))
	render.Render(w, r, res)
}

func (c *RoomController) Create(w http.ResponseWriter, r *http.Request) {
	userId := middleware.UserIdFromContext(r.Context())
	logger := log.WithField("user_id", userId)

	userRepo := c.app.GetRepositories().User
	roomRepo := c.app.GetRepositories().Room
	roleRepo := c.app.GetRepositories().Role

	username, err := userRepo.GetUsername(userId)

	if err != nil {
		logger.WithError(err).Error("Failed to get username of user")
		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	room, err := roomRepo.GetByOwnerId(userId)

	if err != nil {
		log.WithError(err).Error("Failed to get room of user")
		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	if room == nil {
		room = &model.Room{
			Name: fmt.Sprintf("%s's room", username),
			OwnerId: userId,
			Private: false,
		}

		err = roomRepo.Create(room)

		if err != nil {
			logger.WithError(err).Error("Failed to create a room")
			render.Render(w, r, apiResource.ErrInternalError)
			return
		}

		err = roleRepo.Add(&model.UserRole{
			UserId: userId,
			RoomId: room.Id,
			RoleId: role.HOST,
		})

		if err != nil {
			logger.
				WithField("room_id", room.Id).
				WithError(err).
				Error("Failed to add host for a newly created room")

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
		log.
			WithField("room_id", roomId).
			WithError(err).
			Error("Failed to get room")

		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	if room == nil {
		render.Render(w, r, apiResource.ErrNotFound)
		return
	}

	userId := middleware.UserIdFromContext(ctx)

	roleRepo := c.app.GetRepositories().Role
	userRoles, err := roleRepo.Get(userId, room.Id)

	if err != nil {
		log.
			WithFields(log.Fields{
				"user_id": userId,
				"room_id": room.Id,
			}).
			WithError(err).
			Error("Failed to get user roles")

		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	roles := model.BuildRoleManager(userRoles)

	if !roles.HasPermission(permission.MANAGE_ROOM) {
		render.Render(w, r, apiResource.ErrForbidden)
		return
	}

	room.Name = data.Name
	room.Private = data.Private

	err = roomRepo.UpdateInfo(room)

	if err != nil {
		log.
			WithField("room_id", room.Id).
			WithError(err).
			Error("Failed to update room")

		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	cacheKey := fmt.Sprintf(constant.RoomCacheFmt, roomId)
	err = c.app.GetCache().Delete(ctx, cacheKey)

	if err != nil {
		log.
			WithField("room_id", roomId).
			WithError(err).
			Error("Failed to delete room cache")

		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	message := resource.ServerMessage{
		Data: resource.BuildPacket(opcode.UpdateRoom, data),
	}

	bytes, err := msgpack.Marshal(message)

	if err != nil {
		log.
			WithField("message", message).
			WithError(err).
			Error("Failed to serialize message")

		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	roomKey := fmt.Sprintf(constant.RoomFmt, roomId)

	rdb := c.app.GetRedis()
	err = rdb.Publish(ctx, roomKey, bytes).Err()

	if err != nil {
		log.
			WithField("message", message).
			WithError(err).
			Error("Failed to publish message")

		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	render.Render(w, r, apiResource.NewResponse(200))
}

func (c *RoomController) SendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := middleware.UserIdFromContext(ctx)
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
		Author: userId,
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
		log.
			WithField("message", message).
			WithError(err).
			Error("Failed to serialize message")

		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	rdb := c.app.GetRedis()
	err = rdb.Publish(ctx, roomKey, bytes).Err()

	if err != nil {
		log.
			WithField("message", message).
			WithError(err).
			Error("Failed to publish message")

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
			log.
				WithFields(log.Fields{
					"item_id": after,
					"room_id": roomId,
				}).
				WithError(err).
				Error("Failed to get index of queue item")
		} else {
			start += 1 // start 1 element after the specified one
		}
	}

	ids, err := rdb.LRange(ctx, queueKey, start, start + limit - 1).Result()

	if err != nil {
		log.WithError(err).Error("Failed to fetch queue")
		render.Render(w, r, apiResource.ErrInternalError)
		return
	}

	numItems := len(ids)
	queueItemsKey := fmt.Sprintf(constant.RoomQueueItemsFmt, roomId)

	var items []resource.MediaItem

	if numItems > 0 {
		rawItems, err := rdb.HMGet(ctx, queueItemsKey, ids...).Result()

		if err != nil {
			log.WithError(err).Error("Failed to fetch queue items")
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
					log.WithError(err).Error("Failed to parse queue item")
				}
			} else {
				log.WithField("item", rawItem).Warn("Deformed queue item")
			}
		}
	} else {
		items = []resource.MediaItem{}
	}

	res := apiResource.NewQueueResponse(items)
	render.Render(w, r, res)
}