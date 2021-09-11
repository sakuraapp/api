package models

import (
	"gopkg.in/guregu/null.v4"
)

type User struct {
	Id int64
	Username string
	Avatar null.String
	Provider string
	AccessToken null.String
	RefreshToken null.String
	ExternalUserID null.String
	Discriminator null.String `pg:"-"`
}