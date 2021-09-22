package repositories

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/sakuraapp/shared/models"
	"strconv"
)

const MinDiscriminator = "0001"
const MaxDiscriminator = "9999"

type DiscriminatorRepository struct {
	db *pg.DB
}

func (r *DiscriminatorRepository) FindFreeOne(name string) (*string, error) {
	var discrim string

	prev := new(models.Discriminator)
	err := r.db.Model(prev).
		Column("value").
		Where("name = ?", name).
		Order("id DESC").
		First()

	if err == pg.ErrNoRows {
		discrim = MinDiscriminator
		return &discrim, nil
	}

	if err != nil {
		return nil, err
	}

	if prev.Value != MaxDiscriminator {
		intDiscrim, err := strconv.Atoi(prev.Value)

		if err != nil {
			return nil, err
		}

		discrim = fmt.Sprintf("%04d", intDiscrim + 1)

		return &discrim, err
	}

	available := new(models.Discriminator)
	err = r.db.Model(available).
		Column("value").
		Where("owner_id = NULL").
		First()

	if err == pg.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &available.Value, nil
}

func (r *DiscriminatorRepository) Create(discrim *models.Discriminator) error {
	_, err := r.db.Model(discrim).Insert()

	return err
}