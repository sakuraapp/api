package models

type Discriminator struct {
	Id int64
	Name string
	Value string
	OwnerId int64
	Owner *User
}