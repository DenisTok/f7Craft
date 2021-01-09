package src

import (
	"encoding/json"
	"github.com/DenisTok/f7Craft/src/services"
)

type Jsoner interface {
	MetamaskSign(pKey string) ([]byte, error)
	CheckSign(pKey, sign string) ([]byte, error)
	CheckMinecraftSignAndGiveAccess(pKey, sign string) ([]byte, error)

	// from server
	ReqAccess(name string) error
	CheckAccess(name string) error

	// under jwt
	UserProfile(pKey string) ([]byte, error)
	SetMinecraftName(pKey, name string) error
}

type jsoner struct {
	userService     services.UserService
	sessionsService services.SessionsService
	serverService   services.ServerService
}

func NewJsoner(userService services.UserService, sessionsService services.SessionsService, serverService services.ServerService) Jsoner {
	return &jsoner{
		userService:     userService,
		sessionsService: sessionsService,
		serverService:   serverService,
	}
}

func (j *jsoner) MetamaskSign(pKey string) ([]byte, error) {
	res, err := j.userService.MetamaskSign(pKey)
	if err != nil {
		return nil, err
	}

	return json.Marshal(res)
}

func (j *jsoner) CheckSign(pKey, sign string) ([]byte, error) {
	err := j.userService.CheckSign(pKey, sign)
	if err != nil {
		return nil, err
	}

	token, err := j.sessionsService.NewSession(pKey)
	if err != nil {
		return nil, err
	}

	return json.Marshal(token)
}

func (j *jsoner) UserProfile(pKey string) ([]byte, error) {
	profileInfo, err := j.userService.ProfileInfo(pKey)
	if err != nil {
		return nil, err
	}

	return json.Marshal(profileInfo)
}

func (j *jsoner) SetMinecraftName(pKey, name string) error {
	err := j.userService.ChangeMinecraftName(pKey, name)
	if err != nil {
		return err
	}

	return err
}

func (j *jsoner) ReqAccess(name string) error {
	j.serverService.AddToQuarry(name)
	return nil
}

func (j *jsoner) CheckAccess(name string) error {
	err := j.serverService.CheckAccess(name)
	if err != nil {
		return err
	}

	return nil
}

func (j *jsoner) CheckMinecraftSignAndGiveAccess(pKey, sign string) ([]byte, error) {
	name, err := j.userService.CheckMinecraftSign(pKey, sign)
	if err != nil {
		return nil, err
	}

	err = j.serverService.GiveAccess(name)
	if err != nil {
		return nil, err
	}

	return nil, err
}
