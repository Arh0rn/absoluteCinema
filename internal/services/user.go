package services

import (
	"database/sql"
	"errors"
	"github.com/Arh0rn/absoluteCinema/pkg"
	"github.com/Arh0rn/absoluteCinema/pkg/models"
	"github.com/golang-jwt/jwt"
	"github.com/lib/pq"
	"log/slog"
	"strings"
	"time"
)

type UserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByCredentials(email, password string) (*models.User, error)
}

type TokenRepository interface {
	CreateToken(*models.RefreshSession) error
	PopToken(string) (*models.RefreshSession, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
}
type UserServ struct {
	repo      UserRepository
	tokenRepo TokenRepository
	hasher    PasswordHasher

	hmacSecret []byte
	aTokenTTL  time.Duration
	rTokenTTL  time.Duration
}

func NewUserServ(ur UserRepository, tr TokenRepository, h PasswordHasher, secret []byte, attl time.Duration, rttl time.Duration) *UserServ {
	return &UserServ{
		repo:       ur,
		tokenRepo:  tr,
		hasher:     h,
		hmacSecret: secret,
		aTokenTTL:  attl,
		rTokenTTL:  rttl,
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

func (u UserServ) SignIn(inputSignIn models.SignInInput) (string, string, error) {
	password, err := u.hasher.Hash(inputSignIn.Password)
	if err != nil {
		slog.Error(err.Error(),
			"architecture level", "service",
		)
		return "", "", err
	}

	user, err := u.repo.GetUserByCredentials(inputSignIn.Email, password)
	if errors.Is(err, sql.ErrNoRows) {
		slog.Error(models.ErrUserNotFound.Error(),
			"architecture level", "service",
			"email", inputSignIn.Email,
		)
		return "", "", models.ErrUserNotFound
	}
	//if errors.Is(err, sql.Er)
	if err != nil {
		slog.Error(err.Error(),
			"architecture level", "service",
		)
		return "", "", err
	}

	at, rt, err := u.GenerateTokens(user.ID)
	if err != nil {
		slog.Error(err.Error(),
			"architecture level", "service",
		)
		return "", "", err
	}

	return at, rt, nil
}

func (u UserServ) GenerateTokens(userId string) (string, string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   userId,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(u.aTokenTTL).Unix(),
	})
	at, err := token.SignedString(u.hmacSecret)
	if err != nil {
		return "", "", err
	}

	rt, err := pkg.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	err = u.tokenRepo.CreateToken(&models.RefreshSession{
		UserID:    userId,
		Token:     rt,
		ExpiresAt: time.Now().Add(u.rTokenTTL),
	})
	if err != nil {
		return "", "", err
	}

	return at, rt, nil
}

func (u UserServ) RefreshTokens(rt string) (string, string, error) {
	session, err := u.tokenRepo.PopToken(rt)
	if err != nil {
		return "", "", err
	}

	if session.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", models.ErrRefreshTokenExpired
	}

	return u.GenerateTokens(session.UserID)
}
