package repository

import (
	"github.com/go-pg/pg/v10"
	"github.com/sakuraapp/shared/model"
)

type RoomRepository struct {
	db *pg.DB
}

func (r *RoomRepository) Get(id model.RoomId) (*model.Room, error) {
	room := new(model.Room)
	err := r.db.Model(room).
		Relation("Owner").
		Where("room.id = ?", id).
		First()

	if err == pg.ErrNoRows {
		err = nil
		room = nil
	}

	return room, err
}

func (r *RoomRepository) GetLatest() ([]model.Room, error) {
	var rooms []model.Room

	err := r.db.Model(&rooms).
		Relation("Owner").
		Where("private = FALSE").
		Order("id DESC").
		Limit(5).
		Select()

	if err == pg.ErrNoRows {
		err = nil
		rooms = []model.Room{}
	}

	return rooms, err
}

func (r *RoomRepository) GetByOwnerId(id model.UserId) (*model.Room, error) {
	room := new(model.Room)
	err := r.db.Model(room).
		Where("owner_id = ?", id).
		First()

	if err == pg.ErrNoRows {
		err = nil
		room = nil
	}

	return room, err
}

func (r *RoomRepository) Create(room *model.Room) error {
	_, err := r.db.QueryOne(room, `INSERT INTO "rooms" ("id", "name", "owner_id", "private") VALUES (DEFAULT, ?, ?, ?) RETURNING "id"`, room.Name, room.OwnerId, room.Private)

	return err
}