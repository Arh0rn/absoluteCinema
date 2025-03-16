package postgres

import (
	"absoluteCinema/pkg/models"
	"database/sql"
	"log/slog"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

func (r *UserRepo) CreateUser(user *models.User) (*models.User, error) {
	_, err := r.DB.Exec("INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)",
		user.ID, user.Username, user.Email, user.Password)
	if err != nil {
		slog.Error("CreateUser error",
			"architecture level", "repository",
			"error", err.Error(),
		)
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := r.DB.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		slog.Error("GetUserByID error",
			"architecture level", "repository",
			"error", err.Error(),
		)

		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetUserByCredentials(email, password string) (*models.User, error) {
	var user models.User
	err := r.DB.QueryRow("SELECT * FROM users WHERE email = $1 AND password = $2", email, password).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		slog.Error("GetUserByUsername error",
			"architecture level", "repository",
			"error", err.Error(),
		)

		return nil, err
	}
	return &user, nil
}
