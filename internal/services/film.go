package services

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Arh0rn/absoluteCinema/pkg"
	"github.com/Arh0rn/absoluteCinema/pkg/models"
	"log/slog"
)

type FilmRepository interface {
	GetAll() ([]*models.Film, error)
	GetByID(id string) (*models.Film, error)
	UpdateByID(id string, filmInput *models.FilmInput) error
	DeleteByID(id string) error
	Create(*models.Film) error
}
type FilmCache interface {
	GetAll(context.Context) ([]*models.Film, error)
	GetByID(ctx context.Context, id string) (*models.Film, error)
	Set(context.Context, *models.Film) error
	SetAll(context.Context, []*models.Film) error
	Update(ctx context.Context, filmInput *models.Film) error
	Delete(context.Context, string) error
}
type FilmServ struct {
	repo  FilmRepository
	cache FilmCache
}

func NewFilmServ(repo FilmRepository, cache FilmCache) *FilmServ {
	return &FilmServ{
		repo:  repo,
		cache: cache,
	}
}

func (f FilmServ) Create(ctx context.Context, filmDto models.FilmInput) (*models.Film, error) {
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

	err := f.repo.Create(&film)
	if err != nil {
		return nil, err
	}

	err = f.cache.Set(ctx, &film)
	if err != nil {
		slog.Warn("Cache set error",
			"architecture level", "service",
			"error", err.Error())
	}

	return &film, err
}

func (f FilmServ) GetAll(ctx context.Context) ([]*models.Film, error) {
	films, err := f.cache.GetAll(ctx)
	if err == nil && len(films) > 0 {
		return films, nil
	}
	if err != nil {
		slog.Warn("Cache get all error",
			"architecture level", "service",
			"error", err.Error())
	}
	films, err = f.repo.GetAll()

	err = f.cache.SetAll(ctx, films)
	if err != nil {
		slog.Warn("Cache set all error",
			"architecture level", "service",
			"error", err.Error())
	}

	return films, err
}

func (f FilmServ) GetByID(ctx context.Context, id string) (*models.Film, error) {
	film, err := f.cache.GetByID(ctx, id)
	if err == nil && film != nil {
		return film, nil
	}

	film, err = f.repo.GetByID(id)
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

	err = f.cache.Set(ctx, film)
	if err != nil {
		slog.Warn("Cache set error",
			"architecture level", "service",
			"error", err.Error())
	}

	return film, nil
}

func (f FilmServ) UpdateByID(ctx context.Context, id string, filmInput models.FilmInput) error {
	err := f.repo.UpdateByID(id, &filmInput)
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
	film, err := f.repo.GetByID(id)
	if err != nil {
		return err
	}

	err = f.cache.Update(ctx, film)
	if err != nil {
		slog.Warn("Cache delete error",
			"architecture level", "service",
			"error", err.Error())
	}
	return nil
}

func (f FilmServ) DeleteByID(ctx context.Context, id string) error {
	err := f.repo.DeleteByID(id)
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

	err = f.cache.Delete(context.Background(), id)
	if err != nil {
		slog.Warn("Cache delete error",
			"architecture level", "service",
			"error", err.Error())
	}

	return nil
}
