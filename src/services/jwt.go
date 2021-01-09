package services

import (
	"errors"
	"github.com/DenisTok/f7Craft/src/config"
	"github.com/DenisTok/f7Craft/src/models"
	"github.com/DenisTok/f7Craft/src/store/users"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type SessionsService interface {
	NewSession(pAddress string) (token *models.Token, err error)
	RefreshSession(refreshToken string) (token *models.Token, err error)
}

type sessionsService struct {
	sessionsStore *users.SessionsStore
}

func (j *sessionsService) NewSession(pAddress string) (token *models.Token, err error) {
	return j.token(pAddress)
}

func (j *sessionsService) RefreshSession(refreshToken string) (token *models.Token, err error) {
	return nil, errors.New("no")
}

func NewJWTService(sessionsStore *users.SessionsStore) SessionsService {
	return &sessionsService{sessionsStore: sessionsStore}
}

func (j *sessionsService) token(pAddress string) (*models.Token, error) {
	hourToken := time.Now().Add(time.Hour).Unix()
	claim := models.TokenClaims{
		Address:        pAddress,
		StandardClaims: jwt.StandardClaims{ExpiresAt: hourToken},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString(config.SecretKey)
	if err != nil {
		return nil, err
	}

	monthToken := time.Now().Add(time.Hour * 24 * 30).Unix()

	claim.ExpiresAt = monthToken

	rToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	rTokenString, err := rToken.SignedString(config.SecretKey)
	if err != nil {
		return nil, err
	}

	return &models.Token{
		AToken:     tokenString,
		ExpiresAt:  time.Unix(hourToken, 0),
		RToken:     rTokenString,
		RExpiresAt: time.Unix(monthToken, 0),
	}, err
}
