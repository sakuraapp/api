package repository

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/sakuraapp/shared/pkg/model"
	"strconv"
)

const MinDiscriminator = "0001"
const MaxDiscriminator = "9999"

type DiscriminatorRepository struct {
	db *pg.DB
}

func (r *DiscriminatorRepository) FindFreeOne(name string) (*model.Discriminator, error) {
	prev := new(model.Discriminator)
	err := r.db.Model(prev).
		Column("value").
		Where("name = ?", name).
		Order("id DESC").
		First()

	if err == pg.ErrNoRows {
		return &model.Discriminator{
			Name:  name,
			Value: MinDiscriminator,
		}, nil
	}

	if err != nil {
		return nil, err
	}

	if prev.Value != MaxDiscriminator {
		intDiscrim, err := strconv.Atoi(prev.Value)

		if err != nil {
			return nil, err
		}

		return &model.Discriminator{
			Name: name,
			Value: fmt.Sprintf("%04d", intDiscrim + 1),
		}, err
	}

	available := new(model.Discriminator)
	err = r.db.Model(available).
		Column("id", "value").
		Where("owner_id IS NULL").
		First()

	if err == pg.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return available, nil
}

func (r *DiscriminatorRepository) Create(discrim *model.Discriminator) error {
	_, err := r.db.Model(discrim).Insert()

	return err
}

func (r *DiscriminatorRepository) UpdateOwnerID(discrim *model.Discriminator) error {
	_, err := r.db.Model(discrim).
		WherePK().
		Where("owner_id IS NULL").
		Set("owner_id = ?owner_id").
		Update()

	return err
}