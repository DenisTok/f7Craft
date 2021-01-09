package middleware

import (
	"encoding/json"
	"github.com/DenisTok/f7Craft/src/config"
	"github.com/DenisTok/f7Craft/src/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

type jwtBody struct {
	Token string `json:"token"`
}

// APIKey Verify Middleware function
func APIAuth(next fasthttprouter.Handle) fasthttprouter.Handle {
	return func(ctx *fasthttp.RequestCtx, p fasthttprouter.Params) {
		var res jwtBody

		body := ctx.PostBody()

		err := json.Unmarshal(body, &res)
		if err != nil {
			log.Error().Err(err).Send()
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}

		token, err := jwt.ParseWithClaims(res.Token, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			return config.SecretKey, nil
		})
		if err != nil {
			log.Error().Err(err).Send()
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			return
		}

		claims, ok := token.Claims.(*models.TokenClaims)
		if ok && token.Valid {
			ctx.SetUserValue(config.XAddressKey, claims.Address)
		} else {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			return
		}

		next(ctx, p)
	}
}
