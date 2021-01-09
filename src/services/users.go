package services

import (
	"bytes"
	"errors"
	"github.com/DenisTok/f7Craft/src/config"
	"github.com/DenisTok/f7Craft/src/models"
	"github.com/DenisTok/f7Craft/src/store/users"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"time"
)

type UserService interface {
	MetamaskSign(pubAddress string) (res *models.RegRes, err error)
	CheckSign(pKey, sign string) error
	CheckMinecraftSign(pKey, sign string) (string, error)

	// under jwt auth
	ProfileInfo(pKey string) (*models.ProfileInfo, error)
	ChangeMinecraftName(pKey, name string) error
}

type userService struct {
	store *users.Store
}

func NewUserService(store *users.Store) (UserService, error) {
	return &userService{store: store}, nil
}

func (us *userService) MetamaskSign(pubAddress string) (*models.RegRes, error) {
	nonce, err := models.RandEmojis(config.EmojiNonceSize)
	if err != nil {
		return nil, err
	}

	newUser := models.ProtoUser{
		PublicKey: pubAddress,
		Nonce:     nonce + " " + config.NonceAdd,
		Role:      models.Role_Guest,
		Created:   time.Now().UnixNano(),
	}

	ok, err := us.store.UserExist(&newUser)
	if err != nil {
		return nil, err
	}

	if ok {
		err = us.store.UpdateNonce(pubAddress, newUser.Nonce)
		if err != nil {
			return nil, err
		}
	} else {
		err = us.store.SaveUser(&newUser)
		if err != nil {
			return nil, err
		}
	}

	return &models.RegRes{
		Nonce: newUser.Nonce,
	}, nil
}

func (us *userService) CheckSign(pKey, sign string) error {
	u, err := us.store.User(pKey)
	if err != nil {
		return err
	}

	hash := crypto.Keccak256(
		crypto.Keccak256([]byte("string"+" "+"sign")),
		crypto.Keccak256([]byte(u.Nonce)),
	)

	err = us.signChecker(pKey, sign, hash)
	if err != nil {
		return err
	}

	return nil
}

func (us *userService) signChecker(pKey, sign string, hash []byte) error {
	fromAddr := common.HexToAddress(pKey)

	sig := hexutil.MustDecode(sign)

	sig[len(sig)-1] -= 27

	sigPublicKey, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return err
	}

	ethAddress := crypto.PubkeyToAddress(*sigPublicKey)

	if !bytes.Equal(fromAddr.Bytes(), ethAddress.Bytes()) {
		return ErrBadSign
	}

	return nil
}

// return name, err
func (us *userService) CheckMinecraftSign(pKey, sign string) (string, error) {
	u, err := us.store.User(pKey)
	if err != nil {
		return "", err
	}

	if u.MinecraftName == "" {
		return "", ErrNoMinecraftLogin
	}

	hash := crypto.Keccak256(
		crypto.Keccak256([]byte("string"+" "+"game_login")),
		crypto.Keccak256([]byte(u.MinecraftName)),
	)

	err = us.signChecker(pKey, sign, hash)
	if err != nil {
		return "", err
	}

	return u.MinecraftName, nil
}

func (us *userService) ProfileInfo(pKey string) (*models.ProfileInfo, error) {
	u, err := us.store.User(pKey)
	if err != nil {
		return nil, err
	}

	return &models.ProfileInfo{
		AccountAddress: u.PublicKey,
		MinecraftName:  u.MinecraftName,
		Created:        time.Unix(0, u.Created).Format(time.RFC3339),
	}, err
}

func (us *userService) ChangeMinecraftName(pKey, name string) error {
	err := us.store.UpdateMinecraftName(pKey, name)
	if err != nil {
		return err
	}

	return err
}

var ErrBadSign = errors.New("bad sign")
var ErrNoMinecraftLogin = errors.New("no minecraft login")
