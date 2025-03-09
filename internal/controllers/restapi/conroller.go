package restapi

import (
	_ "absoluteCinema/docs"
	"absoluteCinema/internal/controllers/restapi/handlers"
	"net/http"
)

type Controller struct {
	SwaggerController handlers.SwaggerController
	FilmController    handlers.FilmController
}

func NewController(fc *handlers.FilmController) *Controller {
	return &Controller{
		SwaggerController: handlers.SwaggerController{},
		FilmController:    *fc,
	}
}

func (c *Controller) InitRouter() http.Handler {
	mainStackMiddleware := createMiddlewareStack(
		LoggingMiddleware,
		// AuthMiddleware,
		// Other middlewares...
	)

	router := http.NewServeMux()

	// Swagger
	router.Handle("GET /swagger/", c.SwaggerController.Swag())

	{ // Films
		router.HandleFunc("GET /films/", c.FilmController.GetFilms)
		router.HandleFunc("GET /films/{id}", c.FilmController.GetFilmByID)
		router.HandleFunc("POST /films/", c.FilmController.AddFilm)
		router.HandleFunc("PATCH /films/{id}", c.FilmController.UpdateFilmByID)
		router.HandleFunc("DELETE /films/{id}", c.FilmController.DeleteFilmByID)
	}

	return mainStackMiddleware(router)
}
