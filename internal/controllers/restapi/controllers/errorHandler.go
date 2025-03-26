package controllers

import (
	"github.com/Arh0rn/absoluteCinema/pkg/models"
	"log/slog"
	"net/http"
)

func HandleError(w http.ResponseWriter, err error, statusCode int) {
	response := models.ResponseError{
		StatusCode: statusCode,
		Error:      err.Error(),
	}
	w.WriteHeader(statusCode)
	_, writeErr := w.Write([]byte(response.String()))
	if writeErr != nil {
		slog.Error(writeErr.Error(),
			"architecture level", "controller")
		http.Error(w, response.String(), statusCode)
	}
}
