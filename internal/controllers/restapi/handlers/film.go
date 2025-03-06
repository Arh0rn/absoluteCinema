package handlers

import (
	"absoluteCinema/internal/models"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type FilmService interface {
	Create(film models.FilmInput) (*models.Film, error)
	GetAll() ([]models.Film, error)
	GetByID(id string) (*models.Film, error)
	UpdateByID(id string, film models.FilmInput) error
	DeleteByID(id string) error
}

type FilmController struct {
	service FilmService
}

func NewFilmHandler(service FilmService) *FilmController {
	return &FilmController{
		service: service,
	}
}

// AddFilm Add new film
//
//	@Summary      Add a new film
//	@Description  This endpoint adds a new film based on the provided JSON data (excluding ID).
//	@Tags         Films
//	@Accept       json
//	@Produce      json
//	@Param        film  body  models.FilmInput  true  "Film data"
//	@Success      201   {object}  models.Film  "Created successfully"
//	@Failure      400   {string}  string       "Invalid request body"
//	@Failure      500   {string}  string       "Internal Server Error"
//	@Router       /films/ [post]
func (c *FilmController) AddFilm(w http.ResponseWriter, r *http.Request) {
	var film models.FilmInput
	if err := json.NewDecoder(r.Body).Decode(&film); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdFilm, err := c.service.Create(film)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdFilm); err != nil {
		log.Println("Error encoding JSON:", err)
	}
}

// GetFilms returns a list of all films
//
//	@Summary      Get all films
//	@Description  Returns a list of all films in the database.
//	@Tags         Films
//	@Produce      json
//	@Success      200   {array}   models.Film  "Successful response"
//	@Failure      500   {string}  string       "Internal Server Error"
//	@Router       /films/ [get]
func (c *FilmController) GetFilms(w http.ResponseWriter, r *http.Request) {
	films, err := c.service.GetAll()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(films); err != nil {
		log.Println("Error encoding JSON:", err)
	}
}

// GetFilmByID returns a film by ID
//
//	@Summary      Get a film by ID
//	@Description  Returns a single film based on the provided ID.
//	@Tags         Films
//	@Produce      json
//	@Param        id  path  string  true  "Film ID"
//	@Success      200  {object}  models.Film  "Film found"
//	@Failure      400  {string}  string       "Bad Request: ID is required"
//	@Failure      404  {string}  string       "Film not found"
//	@Failure      500  {string}  string       "Internal Server Error"
//	@Router       /films/{id} [get]
func (c *FilmController) GetFilmByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		log.Println("ID is empty")
		http.Error(w, "Bad Request: ID is required", http.StatusBadRequest)
		return
	}

	film, err := c.service.GetByID(id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "Film not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(film); err != nil {
		log.Println("Error encoding JSON:", err)
	}
}

// UpdateFilmByID updates a film by ID
//
//	@Summary      Update a film by ID
//	@Description  Updates an existing film in the database based on the provided ID.
//	@Tags         Films
//	@Accept       json
//	@Param        id    path  string         true  "Film ID"
//	@Param        film  body  models.FilmInput  true  "Updated film data"
//	@Success      204   "Film updated successfully"
//	@Failure      400   {string}  string  "Bad Request: ID is required or invalid JSON"
//	@Failure      500   {string}  string  "Internal Server Error"
//	@Router       /films/{id} [patch]
func (c *FilmController) UpdateFilmByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		log.Println("ID is empty")
		http.Error(w, "Bad Request: ID is required", http.StatusBadRequest)
		return
	}

	var film models.FilmInput
	if err := json.NewDecoder(r.Body).Decode(&film); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.service.UpdateByID(id, film); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteFilmByID deletes a film by ID
//
//	@Summary      Delete a film by ID
//	@Description  Removes a film from the database based on the provided ID.
//	@Tags         Films
//	@Param        id  path  string  true  "Film ID"
//	@Success      204  "Film deleted successfully"
//	@Failure      400  {string}  string  "Bad Request: ID is required"
//	@Failure      500  {string}  string  "Internal Server Error"
//	@Router       /films/{id} [delete]
func (c *FilmController) DeleteFilmByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		log.Println("ID is empty")
		http.Error(w, "Bad Request: ID is required", http.StatusBadRequest)
		return
	}

	if err := c.service.DeleteByID(id); err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
