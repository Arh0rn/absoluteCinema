package app

import (
	"absoluteCinema/internal/controllers/restapi"
	"absoluteCinema/internal/controllers/restapi/handlers"
	"absoluteCinema/internal/repository/postgres"
	"absoluteCinema/internal/services"
	"absoluteCinema/pkg"
	"absoluteCinema/pkg/configParser"
	"absoluteCinema/pkg/database"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

const (
	ConConfigPath = "configs/ConnectionConfig.yaml"
)

type App struct {
	Logger         *slog.Logger
	ConConfig      *configParser.ConnectionConfig
	DB             *sql.DB
	FilmRepository *postgres.FilmRepo
	FilmService    *services.FilmServ
	FilmController *handlers.FilmController
	Router         *restapi.Controller
	Server         *http.Server
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
	ConConf, err := configParser.ParseConnectionConfig(ConConfigPath)
	if err != nil {
		return nil, err
	}
	return ConConf, nil
}

func InitDB() (*sql.DB, error) {
	db, err := database.NewPostgresConnection()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func InitFilmRepository(db *sql.DB) *postgres.FilmRepo {
	FilmRepository := postgres.NewRepo(db)
	return FilmRepository
}

func InitFilmService(FilmRepository *postgres.FilmRepo) *services.FilmServ {
	FilmService := services.NewFilmServ(FilmRepository)
	return FilmService
}

func InitFilmController(FilmService *services.FilmServ) *handlers.FilmController {
	FilmController := handlers.NewFilmHandler(FilmService)
	return FilmController
}

func InitRouter(FilmController *handlers.FilmController) *restapi.Controller {
	router := restapi.NewRouter(FilmController)
	return router
}

func InitServer(ConConf *configParser.ConnectionConfig, router *restapi.Controller) *http.Server {
	middlewares := restapi.CreateMiddlewareStack(
		restapi.LoggingMiddleware,
		//Add more middlewares here
	)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":" + strconv.Itoa(ConConf.Port)),
		Handler: middlewares(router.InitRouter()),
	}
	return srv
}

func InitApp() App {
	var app App

	logger := InitLogger()
	if err := LoadEnv(); err != nil {
		slog.Error(err.Error())
		panic(err)
	}
	ConConf, err := InitConnectionConfig()
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
	db, err := InitDB()
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	FilmRepository := InitFilmRepository(db)
	FilmService := InitFilmService(FilmRepository)
	FilmController := InitFilmController(FilmService)
	router := InitRouter(FilmController)

	srv := InitServer(ConConf, router)

	app = App{
		Logger:         logger,
		ConConfig:      ConConf,
		DB:             db,
		FilmRepository: FilmRepository,
		FilmService:    FilmService,
		FilmController: FilmController,
		Router:         router,
		Server:         srv,
	}

	return app
}

func (a *App) Run() {
	a.Logger.Info("Starting server")
	if err := a.Server.ListenAndServe(); err != nil {
		a.Logger.Error("Server error", err)
	}
}
