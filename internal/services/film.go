package services

import (
	"database/sql"
	"errors"
	"github.com/Arh0rn/absoluteCinema/pkg"
	models2 "github.com/Arh0rn/absoluteCinema/pkg/models"
	"log/slog"
)

type FilmRepository interface {
	GetFilms() ([]models2.Film, error)
	GetFilmByID(id string) (*models2.Film, error)
	UpdateFilmByID(id string, film models2.FilmInput) error
	DeleteFilmByID(id string) error
	CreateFilm(film models2.Film) error
}

type FilmServ struct {
	repo FilmRepository
}

func NewFilmServ(repo FilmRepository) *FilmServ {
	return &FilmServ{
		repo: repo,
	}
}

func (f FilmServ) Create(filmDto models2.FilmInput) (*models2.Film, error) {
	film := models2.Film{
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

func (f FilmServ) GetAll() ([]models2.Film, error) {
	return f.repo.GetFilms()
}

func (f FilmServ) GetByID(id string) (*models2.Film, error) {
	film, err := f.repo.GetFilmByID(id)
	if errors.Is(err, sql.ErrNoRows) {
		slog.Error(models2.ErrFilmNotFound.Error(),
			"architecture level", "service",
			"id", id,
		)
		return nil, models2.ErrFilmNotFound
	}
	if err != nil {
		slog.Error(err.Error(),
			"architecture level", "service",
		)
		return nil, err
	}
	return film, nil
}

func (f FilmServ) UpdateByID(id string, film models2.FilmInput) error {
	err := f.repo.UpdateFilmByID(id, film)
	if errors.Is(err, models2.ErrFilmNotFound) {
		slog.Error(models2.ErrFilmNotFound.Error(),
			"architecture level", "service",
			"id", id,
		)
		return models2.ErrFilmNotFound
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
	if errors.Is(err, models2.ErrFilmNotFound) {
		slog.Error(models2.ErrFilmNotFound.Error(),
			"architecture level", "service",
			"id", id,
		)
		return models2.ErrFilmNotFound
	}
	if err != nil {
		slog.Error(err.Error(),
			"architecture level", "service",
		)
		return err
	}
	return nil
}
