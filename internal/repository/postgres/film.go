package postgres

import (
	models2 "absoluteCinema/pkg/models"
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
)

type FilmRepo struct {
	DB *sql.DB
}

func NewFilmRepo(db *sql.DB) *FilmRepo {
	return &FilmRepo{
		DB: db,
	}
}

func (r *FilmRepo) CreateFilm(film models2.Film) error {
	_, err := r.DB.Exec(
		"INSERT INTO films (id, title, description, release_year, country, duration, budget, box_office) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		film.ID, film.Title, film.Description, film.ReleaseYear, film.Country, film.Duration, film.Budget, film.BoxOffice)
	if err != nil {
		slog.Error("CreateFilm error",
			"architecture level", "repository",
			"error", err.Error(),
		)
		return err
	}
	return nil
}

func (r *FilmRepo) GetFilms() ([]models2.Film, error) {

	rows, err := r.DB.Query("SELECT * FROM films")
	if err != nil {
		slog.Error("GetFilms query error",
			"architecture level", "repository",
			"error", err.Error(),
		)
		return nil, err
	}
	defer rows.Close()

	films := make([]models2.Film, 0)
	for rows.Next() {
		var film models2.Film
		err = rows.Scan(
			&film.ID,
			&film.Title,
			&film.Description,
			&film.ReleaseYear,
			&film.Country,
			&film.Duration,
			&film.Budget,
			&film.BoxOffice)
		if err != nil {
			slog.Error("GetFilms scan error",
				"architecture level", "repository",
				"error", err.Error(),
			)

			return nil, err
		}
		films = append(films, film)
	}
	return films, nil
}

func (r *FilmRepo) GetFilmByID(id string) (*models2.Film, error) {
	var film models2.Film
	err := r.DB.QueryRow("SELECT * FROM films WHERE id = $1", id).Scan(&film.ID, &film.Title, &film.Description, &film.ReleaseYear, &film.Country, &film.Duration, &film.Budget, &film.BoxOffice)
	if err != nil {
		slog.Error("GetFilmByID error",
			"architecture level", "repository",
			"error", err.Error(),
		)

		return nil, err
	}
	return &film, nil
}

func (r *FilmRepo) UpdateFilmByID(id string, film models2.FilmInput) error {
	setValues := make([]string, 0)
	args := make([]any, 0)
	argID := 1

	if film.Title != "" {
		setValues = append(setValues, "title = $"+strconv.Itoa(argID))
		args = append(args, film.Title)
		argID++
	}

	if film.Description != "" {
		setValues = append(setValues, "description = $"+strconv.Itoa(argID))
		args = append(args, film.Description)
		argID++
	}

	if film.ReleaseYear != 0 {
		setValues = append(setValues, "release_year = $"+strconv.Itoa(argID))
		args = append(args, film.ReleaseYear)
		argID++
	}

	if film.Country != "" {
		setValues = append(setValues, "country = $"+strconv.Itoa(argID))
		args = append(args, film.Country)
		argID++
	}

	if film.Duration != 0 {
		setValues = append(setValues, "duration = $"+strconv.Itoa(argID))
		args = append(args, film.Duration)
		argID++
	}

	if film.Budget != 0 {
		setValues = append(setValues, "budget = $"+strconv.Itoa(argID))
		args = append(args, film.Budget)
		argID++
	}

	if film.BoxOffice != 0 {
		setValues = append(setValues, "box_office = $"+strconv.Itoa(argID))
		args = append(args, film.BoxOffice)
		argID++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE films SET %s WHERE id = $%d", setQuery, argID)
	args = append(args, id)

	result, err := r.DB.Exec(query, args...)
	if err != nil {
		slog.Error("UpdateFilmByID error",
			"architecture level", "repository",
			"query", query,
			"args", args,
			"error", err.Error(),
		)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("UpdateFilmByID rows affected error",
			"architecture level", "repository",
			"error", err.Error(),
		)
		return err
	}

	if rowsAffected == 0 {
		slog.Warn("No rows affected",
			"architecture level", "repository",
		)
		return models2.ErrFilmNotFound
	}

	return nil
}

func (r *FilmRepo) DeleteFilmByID(id string) error {
	result, err := r.DB.Exec("DELETE FROM films WHERE id = $1", id)
	if err != nil {
		slog.Error("DeleteFilmByID error",
			"architecture level", "repository",
			"error", err.Error(),
		)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("DeleteFilmByID rows affected error",
			"architecture level", "repository",
			"error", err.Error(),
		)
		return err
	}

	if rowsAffected == 0 {
		slog.Warn("No rows affected",
			"architecture level", "repository",
		)
		return models2.ErrFilmNotFound
	}

	return nil
}
