package utils

import (
	"github.com/gorilla/sessions"
	"net/http"
)

var session *sessions.Session

type FakeStore struct {}

func NewFakeStore() *FakeStore {
	store := &FakeStore{}

	if session == nil {
		session = sessions.NewSession(store, "fake_session")
	}

	return store
}

func (s *FakeStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return session, nil
}

func (s *FakeStore) New(r *http.Request, name string) (*sessions.Session, error) {
	return session, nil
}

func (s *FakeStore) Save(r *http.Request, w http.ResponseWriter, sess *sessions.Session) error  {
	return nil
}