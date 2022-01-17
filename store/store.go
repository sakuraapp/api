package store

import "mime/multipart"

type Service interface {
	Upload(key string, file multipart.File) (location string, err error)
	Delete(key string) error
	ResolveURL(key string) string
}