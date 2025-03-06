package services

import (
	"absoluteCinema/internal/models"
	"absoluteCinema/pkg"
)

type FilmRepository interface {
	GetFilms() ([]models.Film, error)
	GetFilmByID(id string) (*models.Film, error)
	UpdateFilmByID(id string, film models.FilmInput) error
	DeleteFilmByID(id string) error
	CreateFilm(film models.Film) error
}

type FilmServ struct {
	repo FilmRepository
}

func NewFilmServ(repo FilmRepository) *FilmServ {
	return &FilmServ{
		repo: repo,
	}
}

func (f FilmServ) Create(filmDto models.FilmInput) (*models.Film, error) {
	film := models.Film{
		ID:          pkg.GenerateUUID(),
		Title:       filmDto.Title,
		Description: filmDto.Description,
		ReleaseYear: filmDto.ReleaseYear,
		Country:     filmDto.Country,
		Duration:    filmDto.Duration,
		Budget:      filmDto.Budget,
		BoxOffice:   filmDto.BoxOffice,
	}

	err := f.repo.CreateFilm(film)
	if err != nil {
		return nil, err
	}

	return &film, err
}

func (f FilmServ) GetAll() ([]models.Film, error) {
	return f.repo.GetFilms()
}

func (f FilmServ) GetByID(id string) (*models.Film, error) {
	return f.repo.GetFilmByID(id)
}

func (f FilmServ) UpdateByID(id string, film models.FilmInput) error {
	return f.repo.UpdateFilmByID(id, film)
}

func (f FilmServ) DeleteByID(id string) error {
	return f.repo.DeleteFilmByID(id)
}
