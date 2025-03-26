package app

import (
	"database/sql"
	"fmt"
	"github.com/Arh0rn/absoluteCinema/internal/controllers/restapi"
	"github.com/Arh0rn/absoluteCinema/internal/controllers/restapi/controllers"
	"github.com/Arh0rn/absoluteCinema/internal/repository/postgres"
	"github.com/Arh0rn/absoluteCinema/internal/services"
	"github.com/Arh0rn/absoluteCinema/pkg"
	"github.com/Arh0rn/absoluteCinema/pkg/configParser"
	"github.com/Arh0rn/absoluteCinema/pkg/database"
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

func InitTokenRepository(db *sql.DB) *postgres.TokenRepo {
	tokenRepository := postgres.NewTokenRepo(db)
	return tokenRepository
}

func InitUserRepository(db *sql.DB) *postgres.UserRepo {
	userRepository := postgres.NewUserRepo(db)
	return userRepository
}

func InitUserService(ur *postgres.UserRepo, tr *postgres.TokenRepo, h *pkg.Hasher, s []byte, attl time.Duration, rttl time.Duration) *services.UserServ {
	userService := services.NewUserServ(ur, tr, h, s, attl, rttl)
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
