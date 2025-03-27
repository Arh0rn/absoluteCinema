package app

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Arh0rn/absoluteCinema/pkg/configParser"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	ConConfigPath = "configs/ConnectionConfig.yaml"
)

type App struct {
	Logger    *slog.Logger
	ConConfig *configParser.ConnectionConfig
	Server    *http.Server
	DB        *sql.DB
	//FilmRepository *postgres.FilmRepo
	//FilmService    *services.FilmServ
	//FilmController *handlers.FilmController
	//router         http.Handler
}

func InitApp() App {
	var app App

	log := InitLogger()
	log.Info("Init App")
	err := LoadEnv()
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	conConf, err := InitConnectionConfig()
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
	db, err := InitDB()
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
	cache, err := InitCache()
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
	//User
	hasher := InitHasher(os.Getenv("HASH_SALT"))
	secret := []byte(os.Getenv("JWT_SECRET"))
	attl := time.Duration(conConf.AccessTokenTTL) * time.Minute
	rttl := time.Duration(conConf.RefreshTokenTTL) * time.Minute
	userRepository := InitUserRepository(db)
	tokenRepository := InitTokenRepository(db)
	userService := InitUserService(userRepository, tokenRepository, hasher, []byte(secret), attl, rttl)
	userController := InitUserController(userService)

	//Film
	filmRepository := InitFilmRepository(db)
	filmCache := InitFilmCache(cache)
	filmService := InitFilmService(filmRepository, filmCache)
	filmController := InitFilmController(filmService)

	//Server
	controller := InitHandler(filmController, userController)
	router := controller.InitRouter(conConf)
	srv := InitServer(conConf, router)

	app = App{
		Logger:    log,
		ConConfig: conConf,
		Server:    srv,
		DB:        db,
	}

	return app
}

func (a *App) Run() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		a.Logger.Info("Starting server")
		err := a.Server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.Logger.Error("Server start error",
				"error", err.Error(),
			)
		}
	}()

	<-stop
	a.Logger.Info("Shutting down server")

	shutdownTimeout := time.Duration(a.ConConfig.ShutdownTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := a.Server.Shutdown(ctx); err != nil {
		a.Logger.Error("Server shutdown error",
			"error", err,
		)
	}

	if err := a.DB.Close(); err != nil {
		a.Logger.Error("DB close error",
			"error", err,
		)
	}
}
