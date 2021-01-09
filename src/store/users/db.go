package users

import (
	"context"
	"errors"
	"github.com/DenisTok/f7Craft/src/config"
	"github.com/DenisTok/f7Craft/src/models"
	"github.com/DenisTok/f7Craft/src/store"
	"github.com/dgraph-io/badger/v2"
	"github.com/golang/protobuf/proto"
	"strings"
)

type Store struct {
	DB                   *store.DefStore
	MinecraftLoginsIndex *store.DefStore
}

func NewStore(ctx context.Context) (*Store, error) {
	db, err := store.NewDefStore(ctx, config.UsersDBDir)
	if err != nil {
		return nil, err
	}
	indexdb, err := store.NewDefStore(ctx, config.UsersDBDir+"MinecraftLoginsIndex")
	if err != nil {
		return nil, err
	}

	return &Store{DB: db, MinecraftLoginsIndex: indexdb}, err
}

func (s *Store) UpdateNonce(pKey string, nonce string) error {
	u, err := s.User(pKey)
	if err != nil {
		return err
	}

	u.Nonce = nonce

	err = s.SaveUser(u)
	if err != nil {
		return err
	}

	return err
}

func (s *Store) UserExist(u *models.ProtoUser) (bool, error) {
	err := s.DB.Exist([]byte(u.PublicKey))
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *Store) SaveUser(u *models.ProtoUser) error {
	b, err := proto.Marshal(u)
	if err != nil {
		return err
	}

	err = s.DB.Write([]byte(u.PublicKey), b)
	if err != nil {
		return err
	}

	return err
}

func (s *Store) User(publicKey string) (*models.ProtoUser, error) {
	var u models.ProtoUser
	err := s.DB.Get([]byte(publicKey), func(v []byte) error {
		err := proto.Unmarshal(v, &u)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &u, err
}

func (s *Store) UpdateMinecraftName(pKey, name string) error {
	u, err := s.User(pKey)
	if err != nil {
		return err
	}

	if u.MinecraftName != "" {
		return ErrMinecraftNameExist
	}

	var storedAddress string

	err = s.MinecraftLoginsIndex.Get([]byte(strings.ToLower(name)), func(v []byte) error {
		storedAddress = string(v)
		return nil
	})
	if err != nil {
		if !errors.Is(err, badger.ErrKeyNotFound) {
			return err
		}
		storedAddress = u.PublicKey
	}

	if u.PublicKey != storedAddress {
		return ErrMinecraftNameExist
	}

	u.MinecraftName = name

	err = s.SaveUser(u)
	if err != nil {
		return err
	}

	err = s.MinecraftLoginsIndex.Write([]byte(strings.ToLower(name)), []byte(u.PublicKey))
	if err != nil {
		return err
	}

	return err
}

func (s *Store) DeleteMinecraftName(pKey string) error {
	u, err := s.User(pKey)
	if err != nil {
		return err
	}

	var storedAddress string

	err = s.MinecraftLoginsIndex.Get([]byte(strings.ToLower(u.MinecraftName)), func(v []byte) error {
		storedAddress = string(v)
		return nil
	})
	if err != nil {
		return err
	}

	if u.PublicKey != storedAddress {
		return ErrNoAccess
	}

	u.MinecraftName = ""

	err = s.SaveUser(u)
	if err != nil {
		return err
	}

	var keys [][]byte
	keys = append(keys, []byte(strings.ToLower(u.MinecraftName)))

	err = s.MinecraftLoginsIndex.Delete(keys)
	if err != nil {
		return err
	}

	return err
}

var ErrUserExist = errors.New("user exist")
var ErrMinecraftNameExist = errors.New("already exist")
var ErrNoAccess = errors.New("no access")
