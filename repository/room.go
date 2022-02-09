package repository

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-redis/cache/v8"
	"github.com/sakuraapp/shared/constant"
	"github.com/sakuraapp/shared/model"
)

type RoomRepository struct {
	db *pg.DB
	cache *cache.Cache
}

func (r *RoomRepository) Get(ctx context.Context, id model.RoomId) (*model.Room, error) {
	room := new(model.Room)

	if err := r.cache.Once(&cache.Item{
		Ctx:   ctx,
		Key:   fmt.Sprintf(constant.RoomCacheFmt, id),
		Value: room,
		TTL:   constant.RoomCacheTTL,
		Do: func(item *cache.Item) (interface{}, error) {
			return r.fetch(room, id)
		},
	}); err != nil {
		return nil, err
	}

	return room, nil
}

func (r *RoomRepository) fetch(room *model.Room, id model.RoomId) (*model.Room, error) {
	err := r.db.Model(room).
		Column("room.*").
		Relation("Owner").
		ColumnExpr("discriminator.value AS owner__discriminator").
		Join("LEFT JOIN discriminators AS discriminator ON discriminator.owner_id = ?", pg.Ident("owner.id")).
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
		Column("room.*").
		Relation("Owner").
		ColumnExpr("discriminator.value AS owner__discriminator").
		Join("LEFT JOIN discriminators AS discriminator ON discriminator.owner_id = ?", pg.Ident("owner.id")).
		Where("room.private = FALSE").
		Order("room.id DESC").
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

func (r *RoomRepository) UpdateInfo(room *model.Room) error {
	_, err := r.db.Model(room).
		Set("name = ?name, private = ?private").
		Where("id = ?", room.Id).
		Update()

	return err
}