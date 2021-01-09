package models

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Token struct {
	AToken     string    `json:"access_token"`
	ExpiresAt  time.Time `json:"access_exp"`
	RToken     string    `json:"refresh_token"`
	RExpiresAt time.Time `json:"refresh_exp"`
}

type TokenClaims struct {
	Address string `json:"address"`
	jwt.StandardClaims
}
