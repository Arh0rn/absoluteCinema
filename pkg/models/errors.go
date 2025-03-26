package models

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrFilmNotFound         = errors.New("film not found")
	ErrUserNotFound         = errors.New("user not found")
	ErrValidation           = errors.New("validation error")
	ErrInternalServer       = errors.New("internal server error")
	ErrInvalidRequestBody   = errors.New("invalid request body")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrInvalidToken         = errors.New("invalid token")
	ErrRefreshToken         = errors.New("refresh token error")
	ErrRefreshTokenExpired  = errors.New("refresh token expired")
)

type ResponseError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

func (r ResponseError) String() string {
	jsonData, err := json.Marshal(r)
	if err != nil {
		return fmt.Sprintf(`{"error": "failed to serialize error: %v"}`, err)
	}
	return string(jsonData)
}
