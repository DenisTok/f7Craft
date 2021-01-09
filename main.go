package main

import (
	"context"
	"fmt"
	"github.com/DenisTok/f7Craft/src"
	"github.com/DenisTok/f7Craft/src/config"
	"github.com/DenisTok/f7Craft/src/interfaces/rest/application"
	"github.com/DenisTok/f7Craft/src/services"
	"github.com/DenisTok/f7Craft/src/store"
	"github.com/DenisTok/f7Craft/src/store/users"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	dbs []*store.DefStore
}

func main() {
	var app App
	err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Send()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	// setup exit signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Kill, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go func() {
		select {
		case s := <-c:
			fmt.Printf("signal %v", s)
			cancel()
		case <-ctx.Done():
		}
	}()

	// user db init
	userBD, err := users.NewStore(ctx)
	if err != nil {
		log.Fatal().Err(err).Send()
		return
	}

	app.dbs = append(app.dbs, userBD.DB)

	userService, err := services.NewUserService(userBD)
	if err != nil {
		log.Fatal().Err(err).Send()
		return
	}

	// sessions db init
	sessionDB, err := users.NewSessions(ctx)
	if err != nil {
		log.Fatal().Err(err).Send()
		return
	}

	app.dbs = append(app.dbs, sessionDB.DB)

	jwtService := services.NewJWTService(sessionDB)

	serverService := services.NewServerService()

	uc := src.NewJsoner(userService, jwtService, serverService)

	server := application.New(&application.Options{
		Uc:         uc,
		ServerName: "f7Craft-web",
	})

	go server.Start(ctx, ":8080")

	// gracefully shutdown dbs
	select {
	case <-ctx.Done():
		log.Info().Msg("wait for dbs closing")
		for {
			count := 0
			for _, defStore := range app.dbs {
				if defStore.IsClose() {
					count++
				}
			}
			if count == len(app.dbs) {
				break
			}
			time.Sleep(time.Second)
			log.Info().Msg("closing...")
		}
		log.Info().Msg("closed")
	}
}
