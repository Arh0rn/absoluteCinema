package services

import "absoluteCinema/internal/models"

type FilmRepository interface {
	GetFilms() ([]models.Film, error)
	GetFilmByID(id int) (models.Film, error)
	UpdateFilmByID(id int, film models.Film) error
	DeleteFilmByID(id int) error
	CreateFilm(film models.Film) error
}
