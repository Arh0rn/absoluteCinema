package restapi

import (
	_ "absoluteCinema/docs"
	"absoluteCinema/internal/controllers/restapi/controllers"
	"absoluteCinema/pkg/configParser"
	"net/http"
)

type Handler struct {
	SwaggerController controllers.SwaggerController
	FilmController    controllers.FilmController
	UserController    controllers.UserController
}

func NewHandler(fc *controllers.FilmController, uc *controllers.UserController) *Handler {
	return &Handler{
		SwaggerController: controllers.SwaggerController{},
		FilmController:    *fc,
		UserController:    *uc,
	}
}

func (c *Handler) InitRouter(conf *configParser.ConnectionConfig) http.Handler {
	mainStackMiddleware := createMiddlewareStack(
		LoggingMiddleware,
		// AuthMiddleware,
		// Other middlewares...
	)

	router := http.NewServeMux()

	// Swagger
	router.Handle("GET /swagger/", c.SwaggerController.Swag(conf))
	{ // User
		router.HandleFunc("POST /auth/sign-up", c.UserController.SignUp)
		router.HandleFunc("POST /auth/sign-in", c.UserController.SignIn)
	}
	signedRouter := http.NewServeMux()
	{ // Films
		signedRouter.HandleFunc("GET /films/", c.FilmController.GetFilms)
		signedRouter.HandleFunc("GET /films/{id}", c.FilmController.GetFilmByID)
		signedRouter.HandleFunc("POST /films/", c.FilmController.AddFilm)
		signedRouter.HandleFunc("PATCH /films/{id}", c.FilmController.UpdateFilmByID)
		signedRouter.HandleFunc("DELETE /films/{id}", c.FilmController.DeleteFilmByID)
	}
	router.Handle("/", authMiddleware(signedRouter))
	return mainStackMiddleware(router)
}
