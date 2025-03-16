package controllers

import (
	"absoluteCinema/pkg/configParser"
	"absoluteCinema/pkg/models"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"strconv"
)

type SwaggerController struct {
}

var _ models.ResponseError

func (s *SwaggerController) Swag(conf *configParser.ConnectionConfig) http.HandlerFunc {
	url := "http://" + conf.Host + ":" + strconv.Itoa(conf.Port) + "/swagger/doc.json" // http important
	return httpSwagger.Handler(httpSwagger.URL(url))
}
