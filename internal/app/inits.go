package app

import (
	"absoluteCinema/internal/controllers/restapi"
	"absoluteCinema/internal/controllers/restapi/controllers"
	"absoluteCinema/internal/repository/postgres"
	"absoluteCinema/internal/services"
	"absoluteCinema/pkg"
	"absoluteCinema/pkg/configParser"
	"absoluteCinema/pkg/database"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

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

func InitHasher(hash string) *pkg.Hasher {
	hasher := pkg.NewHasher(hash)
	return hasher
}

func InitUserRepository(db *sql.DB) *postgres.UserRepo {
	userRepository := postgres.NewUserRepo(db)
	return userRepository
}

func InitUserService(UserRepository *postgres.UserRepo, hasher *pkg.Hasher, secret []byte, ttl time.Duration) *services.UserServ {
	userService := services.NewUserServ(UserRepository, hasher, secret, ttl)
	return userService
}

func InitUserController(UserService *services.UserServ) *controllers.UserController {
	userController := controllers.NewUserController(UserService)
	return userController
}

func InitFilmRepository(db *sql.DB) *postgres.FilmRepo {
	filmRepository := postgres.NewFilmRepo(db)
	return filmRepository
}

func InitFilmService(FilmRepository *postgres.FilmRepo) *services.FilmServ {
	filmService := services.NewFilmServ(FilmRepository)
	return filmService
}

func InitFilmController(FilmService *services.FilmServ) *controllers.FilmController {
	filmController := controllers.NewFilmController(FilmService)
	return filmController
}

func InitHandler(FilmController *controllers.FilmController, UserController *controllers.UserController) *restapi.Handler {
	controller := restapi.NewHandler(FilmController, UserController)
	return controller
}

func InitServer(ConConf *configParser.ConnectionConfig, router http.Handler) *http.Server {

	srv := &http.Server{
		Addr:    fmt.Sprintf(":" + strconv.Itoa(ConConf.Port)),
		Handler: router,
	}
	return srv
}
