package services

import (
	"absoluteCinema/pkg"
	"absoluteCinema/pkg/models"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"log/slog"
	"strings"
	"time"
)

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByCredentials(email, password string) (*models.User, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}
type UserServ struct {
	repo   UserRepository
	hasher PasswordHasher

	hmacSecret []byte
	tokenTTL   time.Duration
}

func NewUserServ(repo UserRepository, hasher PasswordHasher, hmacSecret []byte, tokenTTL time.Duration) *UserServ {
	return &UserServ{
		repo:       repo,
		hasher:     hasher,
		hmacSecret: hmacSecret,
		tokenTTL:   tokenTTL,
	}
}

func (u UserServ) SignUp(signUpInput models.SignUpInput) (*models.User, error) {
	passwordHash, err := u.hasher.Hash(signUpInput.Password)
	if err != nil {
		slog.Error(err.Error(),
			"architecture level", "service",
		)
		return nil, err
	}

	user := &models.User{
		ID:       pkg.GenerateUUID(),
		Username: signUpInput.Username,
		Email:    signUpInput.Email,
		Password: passwordHash,
	}

	createdUser, err := u.repo.CreateUser(user)

	if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
		if strings.Contains(pgErr.Constraint, "users_email_key") {
			slog.Error(models.ErrUserAlreadyExists.Error(),
				"architecture level", "service",
				"email", signUpInput.Email,
			)
			return nil, models.ErrUserAlreadyExists
		}
		if strings.Contains(pgErr.Constraint, "users_username_key") {
			slog.Error(models.ErrUsernameAlreadyTaken.Error(),
				"architecture level", "service",
				"username", signUpInput.Username,
			)
			return nil, models.ErrUsernameAlreadyTaken
		}
		slog.Error(models.ErrUserAlreadyExists.Error(),
			"architecture level", "service",
			"email", signUpInput.Email,
		)

	}
	if err != nil {
		slog.Error(err.Error(),
			"architecture level", "service",
		)
		return nil, err
	}

	return createdUser, nil
}

func (u UserServ) SignIn(inputSignIn models.SignInInput) (string, error) {
	password, err := u.hasher.Hash(inputSignIn.Password)
	if err != nil {
		slog.Error(err.Error(),
			"architecture level", "service",
		)
		return "", err
	}

	user, err := u.repo.GetUserByCredentials(inputSignIn.Email, password)
	if errors.Is(err, sql.ErrNoRows) {
		slog.Error(models.ErrUserNotFound.Error(),
			"architecture level", "service",
			"email", inputSignIn.Email,
		)
		return "", models.ErrUserNotFound
	}
	//if errors.Is(err, sql.Er)
	if err != nil {
		slog.Error(err.Error(),
			"architecture level", "service",
		)
		return "", err
	}

	token, err := pkg.GenerateToken(user, u.hmacSecret, u.tokenTTL)
	if err != nil {
		slog.Error(err.Error(),
			"architecture level", "service",
		)
		return "", err
	}

	return token, nil
}

func (u UserServ) ParseToken(token string) (string, error) {
	id, err := pkg.ParseToken(token, u.hmacSecret)
	if err != nil {
		slog.Error(err.Error(),
			"architecture level", "service",
		)
		return "", err
	}

	return id, nil
}
