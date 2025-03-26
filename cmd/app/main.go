package main

import (
	"github.com/Arh0rn/absoluteCinema/internal/app"
)

// @title absoluteCinema
// @version 1.5
// @description This is a sample server Films data  server.
// @contact.email amir.kurmanbekov@gmail.com
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	a := app.InitApp()
	a.Run()
}
