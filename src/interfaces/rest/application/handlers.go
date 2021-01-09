package application

import (
	"encoding/json"
	"errors"
	"github.com/DenisTok/f7Craft/src/config"
	"github.com/DenisTok/f7Craft/src/models"
	"github.com/DenisTok/f7Craft/src/services"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"
)

// HealthHandler check is alive
func (app *Application) HealthHandler(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (app *Application) GenNonce(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	pKey := string(ctx.QueryArgs().Peek("public_key"))
	if pKey == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	res, err := app.uc.MetamaskSign(pKey)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		log.Error().Err(err).Send()
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	_, _ = ctx.Write(res)
}

type SignReq struct {
	Sign string `json:"sign"`
}

func (app *Application) CheckSign(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	pKey := string(ctx.QueryArgs().Peek("public_key"))
	if pKey == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	var req SignReq

	err := json.Unmarshal(ctx.PostBody(), &req)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	res, err := app.uc.CheckSign(pKey, req.Sign)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		log.Error().Err(err).Send()
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	_, _ = ctx.Write(res)
}

func (app *Application) UserProfile(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	address, ok := ctx.UserValue(config.XAddressKey).(string)
	if !ok {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		log.Error().Msg("incorrect interface when UserProfile")
		return
	}

	res, err := app.uc.UserProfile(address)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		log.Error().Err(err).Send()
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	_, _ = ctx.Write(res)
}

func (app *Application) SetMinecraftName(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	address, ok := ctx.UserValue(config.XAddressKey).(string)
	if !ok {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		log.Error().Msg("incorrect interface when SetMinecraftName")
		return
	}

	var req models.MinecraftNameChangeReq

	err := json.Unmarshal(ctx.PostBody(), &req)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		log.Error().Err(err).Send()
		return
	}

	err = app.uc.SetMinecraftName(address, req.Name)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		log.Error().Err(err).Send()
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (app *Application) CheckMinecraftSign(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	pKey := string(ctx.QueryArgs().Peek("public_key"))
	if pKey == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	var req SignReq

	err := json.Unmarshal(ctx.PostBody(), &req)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	res, err := app.uc.CheckMinecraftSignAndGiveAccess(pKey, req.Sign)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		log.Error().Err(err).Send()
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	_, _ = ctx.Write(res)
}

func (app *Application) ReqAccessToServer(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	name := string(ctx.QueryArgs().Peek("login"))
	if name == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	err := app.uc.ReqAccess(name)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		log.Error().Err(err).Send()
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}

func (app *Application) CheckAccessToServer(ctx *fasthttp.RequestCtx, _ fasthttprouter.Params) {
	name := string(ctx.QueryArgs().Peek("login"))
	if name == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	err := app.uc.CheckAccess(name)
	if err != nil {
		if errors.Is(err, services.ErrWait) {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			return
		}
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		log.Error().Err(err).Send()
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
}
