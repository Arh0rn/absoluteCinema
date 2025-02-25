package main

import (
	"absoluteCinema/pkg"
	"log"
)

func main() {
	if err := pkg.LoadEnv(); err != nil {
		log.Fatal("Env load error", err)
	}

	// filmRepo := postgres.NewRepo(db)
	// filmService := services.NewFilmService(filmRepo)
	// filmHandler := handlers.NewFilmHandler(filmService)
	// filmHandler.InitRoutes(router)
	// log.Fatal(http.ListenAndServe(":8080", router))
}
