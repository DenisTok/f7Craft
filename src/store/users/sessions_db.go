package users

import (
	"context"
	"errors"
	"github.com/DenisTok/f7Craft/src/config"
	"github.com/DenisTok/f7Craft/src/models"
	"github.com/DenisTok/f7Craft/src/store"
	"github.com/dgraph-io/badger/v2"

	"github.com/golang/protobuf/proto"
)

type SessionsStore struct {
	DB *store.DefStore
}

func NewSessions(ctx context.Context) (*SessionsStore, error) {
	db, err := store.NewDefStore(ctx, config.SessionsUsersDBDir)
	if err != nil {
		return nil, err
	}

	return &SessionsStore{DB: db}, err
}

func (ss *SessionsStore) SaveSession(s *models.ProtoSession) error {
	sessions, err := ss.GetSessions(s.PublicKey)
	if err != nil {
		if !errors.Is(err, badger.ErrKeyNotFound) {
			return err
		}
		sessions = &models.ProtoSessions{}
	}

	sessions.ProtoSessions = append(sessions.ProtoSessions, s)

	b, err := proto.Marshal(sessions)
	if err != nil {
		return err
	}

	err = ss.DB.Write([]byte(s.PublicKey), b)
	if err != nil {
		return err
	}

	return err
}

func (ss *SessionsStore) GetSessions(pKey string) (*models.ProtoSessions, error) {
	var sessions models.ProtoSessions

	err := ss.DB.Get([]byte(pKey), func(v []byte) error {
		err := proto.Unmarshal(v, &sessions)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &sessions, err
}
