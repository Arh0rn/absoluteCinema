package postgres

import (
	"absoluteCinema/internal/models"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

type Repo struct {
	DB *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{
		DB: db,
	}
}

//	CreateFilm(film models.Film) error
//	GetFilms() ([]models.Film, error)
//	GetFilmByID(id int) (models.Film, error)
//	UpdateFilmByID(id int, film models.Film) error
//	DeleteFilmByID(id int) error

//type Film struct {
//	ID          string   `json:"id"`
//	Title       string   `json:"title"`
//	Description string   `json:"description"`
//	ReleaseYear int      `json:"release_year"`
//	Country     string   `json:"country"`
//	Duration    int      `json:"duration"`
//	Budget      int      `json:"budget"`
//	BoxOffice   int      `json:"box_office"`
//}

func (r *Repo) CreateFilm(film models.Film) error {
	_, err := r.DB.Exec(
		"INSERT INTO films (title, description, release_year, country, duration, budget, box_office) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		film.Title, film.Description, film.ReleaseYear, film.Country, film.Duration, film.Budget, film.BoxOffice)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) GetFilms() ([]models.Film, error) {
	rows, err := r.DB.Query("SELECT * FROM films")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	films := make([]models.Film, 0)
	for rows.Next() {
		var film models.Film
		err = rows.Scan(&film.ID, &film.Title, &film.Description, &film.ReleaseYear, &film.Country, &film.Duration, &film.Budget, &film.BoxOffice)
		if err != nil {
			return nil, err
		}
		films = append(films, film)
	}
	return films, nil
}

func (r *Repo) GetFilmByID(id int) (models.Film, error) {
	var film models.Film
	err := r.DB.QueryRow("SELECT * FROM films WHERE id = $1", id).Scan(&film.ID, &film.Title, &film.Description, &film.ReleaseYear, &film.Country, &film.Duration, &film.Budget, &film.BoxOffice)
	if err != nil {
		return models.Film{}, err
	}
	return film, nil
}

func (r *Repo) UpdateFilmByID(id int, film models.Film) error {
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

	_, err := r.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) DeleteFilmByID(id int) error {
	_, err := r.DB.Exec("DELETE FROM films WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
