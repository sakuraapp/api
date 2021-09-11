package models

type Room struct {
	Id int64
	Name string
	OwnerId int64 ``
	Owner *User `pg:"rel:has-one"`
	Private bool
}