package application

import (
	"context"
	"github.com/DenisTok/f7Craft/src/interfaces/rest/middleware"
	"os"
	"time"

	"github.com/lab259/cors"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttprouter"

	"github.com/DenisTok/f7Craft/src"
)

// Application struct
type Application struct {
	r    *fasthttprouter.Router
	serv *fasthttp.Server
	uc   src.Jsoner
}

// Options struct
type Options struct {
	Uc         src.Jsoner
	ServerName string
}

// New return new application
func New(opt *Options) *Application {
	router := fasthttprouter.New()

	return &Application{
		r: router,
		serv: &fasthttp.Server{
			ReadTimeout:  time.Second * 5,
			IdleTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 5,
			Handler:      cors.AllowAll().Handler(router.Handler),
		},
		uc: opt.Uc,
	}
}

func (app *Application) Start(ctx context.Context, httpPort string) {

	app.r.GET("/alive", app.HealthHandler)
	app.r.POST("/api/v1/users/nonce", app.GenNonce)
	app.r.POST("/api/v1/users/sign", app.CheckSign)
	app.r.POST("/api/v1/mc/sign", app.CheckMinecraftSign)

	app.r.GET("/xserver/access/req", app.ReqAccessToServer)
	app.r.GET("/xserver/access/check", app.CheckAccessToServer)

	app.r.POST("/api/v1/user/profile", middleware.APIAuth(app.UserProfile))
	app.r.POST("/api/v1/user/mc/name", middleware.APIAuth(app.SetMinecraftName))

	listenErr := make(chan error, 1)

	go func() {
		listenErr <- app.serv.ListenAndServe(httpPort)
	}()

	log.Info().Msg("http ok")

	select {
	case err := <-listenErr:
		if err != nil {
			log.Error().Msgf("listen err %s", err)
			os.Exit(1)
		}
	case <-ctx.Done():
		if err := app.serv.Shutdown(); err != nil {
			log.Error().Msgf("shutdown err %s", err)
			os.Exit(1)
		}
	}
}
