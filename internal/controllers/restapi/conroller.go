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

func NewRouter(fc *handlers.FilmController) *Controller {
	return &Controller{
		SwaggerController: handlers.SwaggerController{},
		FilmController:    *fc,
	}
}

func (c *Controller) InitRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /swagger/", c.SwaggerController.Swag())
	mux.HandleFunc("GET /films/", c.FilmController.GetFilms)
	mux.HandleFunc("GET /films/{id}", c.FilmController.GetFilmByID)
	mux.HandleFunc("POST /films/", c.FilmController.AddFilm)
	mux.HandleFunc("PATCH /films/{id}", c.FilmController.UpdateFilmByID)
	mux.HandleFunc("DELETE /films/{id}", c.FilmController.DeleteFilmByID)

	return AddMiddlewares(mux)
}
