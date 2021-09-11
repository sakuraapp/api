package repositories

import (
	"github.com/go-pg/pg/v10"
	"github.com/sakuraapp/api/models"
)

type RoomRepository struct {
	db *pg.DB
}

func (r *RoomRepository) Get(id int64) (*models.Room, error) {
	room := new(models.Room)
	err := r.db.Model(&room).
		Relation("Owner").
		Where("id = ?", id).
		First()

	if err == pg.ErrNoRows {
		err = nil
		room = nil
	}

	return room, err
}

func (r *RoomRepository) GetLatest() ([]models.Room, error) {
	var rooms []models.Room

	err := r.db.Model(&rooms).
		Relation("Owner").
		Where("private = FALSE").
		Order("id DESC").
		Limit(5).
		Select()

	if err == pg.ErrNoRows {
		err = nil
		rooms = []models.Room{}
	}

	return rooms, err
}

func (r *RoomRepository) GetByOwnerId(id int64) (*models.Room, error) {
	room := new(models.Room)
	err := r.db.Model(room).
		Where("owner_id = ?", id).
		First()

	if err == pg.ErrNoRows {
		err = nil
		room = nil
	}

	return room, err
}

func (r *RoomRepository) Create(room *models.Room) error {
	_, err := r.db.Model(room).Insert()

	return err
}