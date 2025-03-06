package handlers

import (
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type SwaggerController struct {
}

func (s *SwaggerController) Swag() http.HandlerFunc {
	return httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json"))
}
