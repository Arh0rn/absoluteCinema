package app

import (
	"absoluteCinema/internal/controllers/restapi"
	"absoluteCinema/internal/controllers/restapi/handlers"
	"absoluteCinema/internal/repository/postgres"
	"absoluteCinema/internal/services"
	"absoluteCinema/pkg"
	"absoluteCinema/pkg/configParser"
	"absoluteCinema/pkg/database"
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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

func InitLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	return logger
}

func LoadEnv() error {
	if err := pkg.LoadEnv(); err != nil {
		return err
	}
	return nil
}

func InitConnectionConfig() (*configParser.ConnectionConfig, error) {
	conConf, err := configParser.ParseConnectionConfig(ConConfigPath)
	if err != nil {
		return nil, err
	}
	return conConf, nil
}

func InitDB() (*sql.DB, error) {
	db, err := database.NewPostgresConnection()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitFilmRepository(db *sql.DB) *postgres.FilmRepo {
	filmRepository := postgres.NewRepo(db)
	return filmRepository
}

func InitFilmService(FilmRepository *postgres.FilmRepo) *services.FilmServ {
	filmService := services.NewFilmServ(FilmRepository)
	return filmService
}

func InitFilmController(FilmService *services.FilmServ) *handlers.FilmController {
	filmController := handlers.NewFilmController(FilmService)
	return filmController
}

func InitController(FilmController *handlers.FilmController) *restapi.Controller {
	controller := restapi.NewController(FilmController)
	return controller
}

func InitServer(ConConf *configParser.ConnectionConfig, router http.Handler) *http.Server {

	srv := &http.Server{
		Addr:    fmt.Sprintf(":" + strconv.Itoa(ConConf.Port)),
		Handler: router,
	}
	return srv
}

func InitApp() App {
	var app App

	logger := InitLogger()

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

	filmRepository := InitFilmRepository(db)
	filmService := InitFilmService(filmRepository)
	filmController := InitFilmController(filmService)
	controller := InitController(filmController)
	router := controller.InitRouter()
	srv := InitServer(conConf, router)

	app = App{
		Logger:    logger,
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
