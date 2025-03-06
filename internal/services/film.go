package services

import (
	"absoluteCinema/internal/models"
	"absoluteCinema/pkg"
	"database/sql"
	"errors"
	"log/slog"
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
	film, err := f.repo.GetFilmByID(id)
	if errors.Is(err, sql.ErrNoRows) {
		slog.Error(models.ErrFilmNotFound.Error(),
			"architecture level", "service",
			"id", id,
		)
		return nil, models.ErrFilmNotFound
	}
	if err != nil {
		slog.Error(err.Error(),
			"architecture level", "service",
		)
		return nil, err
	}
	return film, nil
}

func (f FilmServ) UpdateByID(id string, film models.FilmInput) error {
	err := f.repo.UpdateFilmByID(id, film)
	if errors.Is(err, models.ErrFilmNotFound) {
		slog.Error(models.ErrFilmNotFound.Error(),
			"architecture level", "service",
			"id", id,
		)
		return models.ErrFilmNotFound
	}
	if err != nil {
		slog.Error(err.Error(),
			"architecture level", "service",
		)
		return err
	}
	return nil
}

func (f FilmServ) DeleteByID(id string) error {
	err := f.repo.DeleteFilmByID(id)
	if errors.Is(err, models.ErrFilmNotFound) {
		slog.Error(models.ErrFilmNotFound.Error(),
			"architecture level", "service",
			"id", id,
		)
		return models.ErrFilmNotFound
	}
	if err != nil {
		slog.Error(err.Error(),
			"architecture level", "service",
		)
		return err
	}
	return nil
}
