package controllers

import (
	"encoding/json"
	"errors"
	"github.com/Arh0rn/absoluteCinema/pkg/models"
	"log/slog"
	"net/http"
)

type UserService interface {
	SignUp(signUpInput models.SignUpInput) (*models.User, error)
	SignIn(inputSignIn models.SignInInput) (string, string, error)
	RefreshTokens(rt string) (string, string, error)
}

type UserController struct {
	service UserService
}

func NewUserController(service UserService) *UserController {
	return &UserController{
		service: service,
	}
}

// SignUp Create a new user account
//
//	@Summary      Sign up a new user
//	@Description  This endpoint is responsible for registering a new user.
//	@Tags         Auth
//	@Accept       json
//	@Produce      json
//	@Param        user  body  models.SignUpInput  true  "User data"
//	@Success      201   {object}  models.User  "User created successfully"
//	@Failure      400   {object}  models.ResponseError  "Invalid request body"
//	@Failure      422   {object}  models.ResponseError  "Validation error"
//	@Failure      500   {object}  models.ResponseError  "Internal Server Error"
//	@Router       /auth/sign-up [post]
func (c *UserController) SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var signUpInput models.SignUpInput
	if err := json.NewDecoder(r.Body).Decode(&signUpInput); err != nil {
		slog.Error(err.Error(),
			"architecture level", "controller",
		)
		//http.Error(w, models.REInvalidRequestBody.String(), http.StatusBadRequest)
		HandleError(w, models.ErrInvalidRequestBody, http.StatusBadRequest)
		return
	}

	if err := signUpInput.Validate(); err != nil {
		slog.Error(err.Error(),
			"architecture level", "controller",
		)
		//http.Error(w, models.REValidation.String(), http.StatusUnprocessableEntity)
		HandleError(w, models.ErrValidation, http.StatusUnprocessableEntity)
		return
	}

	user, err := c.service.SignUp(signUpInput)
	if errors.Is(err, models.ErrUserAlreadyExists) {
		slog.Error(models.ErrUserAlreadyExists.Error(),
			"architecture level", "controller",
		)
		//http.Error(w, models.REUserAlreadyExists.String(), http.StatusConflict)
		HandleError(w, models.ErrUserAlreadyExists, http.StatusBadRequest)
		return
	}
	if errors.Is(err, models.ErrUsernameAlreadyTaken) {
		slog.Error(models.ErrUsernameAlreadyTaken.Error(),
			"architecture level", "controller",
		)
		//http.Error(w, models.REUsernameAlreadyTaken.String(), http.StatusConflict)
		HandleError(w, models.ErrUsernameAlreadyTaken, http.StatusBadRequest)
		return
	}
	if err != nil {
		slog.Error(err.Error(),
			"architecture level", "controller",
		)
		//http.Error(w, models.REInternalServer.String(), http.StatusInternalServerError)
		HandleError(w, models.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		slog.Error("Error encoding JSON",
			"architecture level", "controller",
			"error", err.Error(),
		)
		return
	}
}

// SignIn Authenticate a user and return a token
//
//	@Summary      Sign in a user
//	@Description  This endpoint allows users to authenticate and receive an access token.
//	@Tags         Auth
//	@Accept       json
//	@Produce      json
//	@Param        user  body  models.SignInInput  true  "User credentials"
//	@Success      200   {object}  map[string]string  "Authentication successful, returns access token"
//	@Failure      400   {object}  models.ResponseError  "Invalid request body"
//	@Failure      401   {object}  models.ResponseError  "Invalid credentials"
//	@Failure      404   {object}  models.ResponseError  "User not found"
//	@Failure      409   {object}  models.ResponseError  "User already exists"
//	@Failure      409   {object}  models.ResponseError  "Username already taken"
//	@Failure      422   {object}  models.ResponseError  "Validation error"
//	@Failure      500   {object}  models.ResponseError  "Internal Server Error"
//	@Router       /auth/sign-in [post]
func (c *UserController) SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var inputSignIn models.SignInInput
	if err := json.NewDecoder(r.Body).Decode(&inputSignIn); err != nil {
		slog.Error(err.Error(),
			"architecture level", "controller",
		)
		//http.Error(w, models.REInvalidRequestBody.String(), http.StatusBadRequest)
		HandleError(w, models.ErrInvalidRequestBody, http.StatusBadRequest)
		return
	}

	if err := inputSignIn.Validate(); err != nil {
		slog.Error(err.Error(),
			"architecture level", "controller",
		)
		//http.Error(w, models.REValidation.String(), http.StatusUnprocessableEntity)
		HandleError(w, models.ErrValidation, http.StatusUnprocessableEntity)
	}
	at, rt, err := c.service.SignIn(inputSignIn)
	if errors.Is(err, models.ErrUserNotFound) {
		slog.Error(models.ErrUserNotFound.Error(),
			"architecture level", "controller",
		)
		//http.Error(w, models.REUserNotFound.String(), http.StatusNotFound)
		HandleError(w, models.ErrUserNotFound, http.StatusBadRequest)
		return
	}
	if err != nil {
		slog.Error(err.Error(),
			"architecture level", "controller",
		)
		//http.Error(w, models.REInvalidRequestBody.String(), http.StatusInternalServerError)
		HandleError(w, models.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Set-Cookie", "refresh-token="+rt+"; HttpOnly")
	err = json.NewEncoder(w).Encode(map[string]string{"token": at})
	if err != nil {
		slog.Error("Error encoding JSON",
			"architecture level", "controller",
			"error", err.Error(),
		)
		HandleError(w, models.ErrInternalServer, http.StatusInternalServerError)
		return
	}
}

func (c *UserController) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh-token")
	if err != nil {
		slog.Error(err.Error(),
			"architecture level", "controller",
			"cookie value", cookie,
		)
		HandleError(w, models.ErrRefreshToken, http.StatusBadRequest)
		return
	}
	//slog.Info("Refresh token", "cookie value", cookie)
	at, rt, err := c.service.RefreshTokens(cookie.Value)
	if errors.Is(err, models.ErrRefreshTokenExpired) {
		slog.Error(err.Error(),
			"architecture level", "controller",
		)
		HandleError(w, models.ErrRefreshTokenExpired, http.StatusUnauthorized)
		return
	}
	if err != nil {
		slog.Error(err.Error(),
			"architecture level", "controller",
		)
		HandleError(w, models.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Set-Cookie", "refresh-token="+rt+"; HttpOnly")
	err = json.NewEncoder(w).Encode(map[string]string{"token": at})
	if err != nil {
		slog.Error("Error encoding JSON",
			"architecture level", "controller",
			"error", err.Error(),
		)
		HandleError(w, models.ErrInternalServer, http.StatusInternalServerError)
	}
}
