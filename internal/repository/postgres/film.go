package postgres

import (
	"database/sql"
	"fmt"
	"github.com/Arh0rn/absoluteCinema/pkg/models"
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

func (r *FilmRepo) Create(film *models.Film) error {
	_, err := r.DB.Exec(
		"INSERT INTO films (id, title, description, release_year, country, duration, budget, box_office) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		film.ID, film.Title, film.Description, film.ReleaseYear, film.Country, film.Duration, film.Budget, film.BoxOffice)
	if err != nil {
		slog.Error("Create error",
			"architecture level", "repository",
			"error", err.Error(),
		)
		return err
	}
	return nil
}

func (r *FilmRepo) GetAll() ([]*models.Film, error) {

	rows, err := r.DB.Query("SELECT * FROM films")
	if err != nil {
		slog.Error("GetAll query error",
			"architecture level", "repository",
			"error", err.Error(),
		)
		return nil, err
	}
	defer rows.Close()

	films := make([]*models.Film, 0)
	for rows.Next() {
		var film models.Film
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
			slog.Error("GetAll scan error",
				"architecture level", "repository",
				"error", err.Error(),
			)

			return nil, err
		}
		films = append(films, &film)
	}
	return films, nil
}

func (r *FilmRepo) GetByID(id string) (*models.Film, error) {
	var film models.Film
	err := r.DB.QueryRow("SELECT * FROM films WHERE id = $1", id).Scan(&film.ID, &film.Title, &film.Description, &film.ReleaseYear, &film.Country, &film.Duration, &film.Budget, &film.BoxOffice)
	if err != nil {
		slog.Error("GetByID error",
			"architecture level", "repository",
			"error", err.Error(),
		)

		return nil, err
	}
	return &film, nil
}

func (r *FilmRepo) UpdateByID(id string, filmInput *models.FilmInput) error {
	setValues := make([]string, 0)
	args := make([]any, 0)
	argID := 1

	if filmInput.Title != "" {
		setValues = append(setValues, "title = $"+strconv.Itoa(argID))
		args = append(args, filmInput.Title)
		argID++
	}

	if filmInput.Description != "" {
		setValues = append(setValues, "description = $"+strconv.Itoa(argID))
		args = append(args, filmInput.Description)
		argID++
	}

	if filmInput.ReleaseYear != 0 {
		setValues = append(setValues, "release_year = $"+strconv.Itoa(argID))
		args = append(args, filmInput.ReleaseYear)
		argID++
	}

	if filmInput.Country != "" {
		setValues = append(setValues, "country = $"+strconv.Itoa(argID))
		args = append(args, filmInput.Country)
		argID++
	}

	if filmInput.Duration != 0 {
		setValues = append(setValues, "duration = $"+strconv.Itoa(argID))
		args = append(args, filmInput.Duration)
		argID++
	}

	if filmInput.Budget != 0 {
		setValues = append(setValues, "budget = $"+strconv.Itoa(argID))
		args = append(args, filmInput.Budget)
		argID++
	}

	if filmInput.BoxOffice != 0 {
		setValues = append(setValues, "box_office = $"+strconv.Itoa(argID))
		args = append(args, filmInput.BoxOffice)
		argID++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE films SET %s WHERE id = $%d", setQuery, argID)
	args = append(args, id)

	result, err := r.DB.Exec(query, args...)
	if err != nil {
		slog.Error("UpdateByID error",
			"architecture level", "repository",
			"query", query,
			"args", args,
			"error", err.Error(),
		)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("UpdateByID rows affected error",
			"architecture level", "repository",
			"error", err.Error(),
		)
		return err
	}

	if rowsAffected == 0 {
		slog.Warn("No rows affected",
			"architecture level", "repository",
		)
		return models.ErrFilmNotFound
	}

	return nil
}

func (r *FilmRepo) DeleteByID(id string) error {
	result, err := r.DB.Exec("DELETE FROM films WHERE id = $1", id)
	if err != nil {
		slog.Error("DeleteByID error",
			"architecture level", "repository",
			"error", err.Error(),
		)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("DeleteByID rows affected error",
			"architecture level", "repository",
			"error", err.Error(),
		)
		return err
	}

	if rowsAffected == 0 {
		slog.Warn("No rows affected",
			"architecture level", "repository",
		)
		return models.ErrFilmNotFound
	}

	return nil
}
