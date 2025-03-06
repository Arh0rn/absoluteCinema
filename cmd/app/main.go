package main

import (
	"absoluteCinema/internal/app"
)

// @title absoluteCinema
// @version 1.0
// @description This is a sample server Films data  server.
// @contact.email amir.kurmanbekov@gmail.com
// @host localhost:8080
// @BasePath /

func main() {
	a := app.InitApp()
	a.Run()
}

//func main() {
//	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
//	slog.SetDefault(logger)
//
//	if err := pkg.LoadEnv(); err != nil {
//		log.Fatal("Env load error", err)
//	}
//	//dbConf, err := configParser.ParseDBConfig("/configs/dbConfig.yaml")
//	//if err != nil {
//	//	log.Fatal("Config parse error", err)
//	//}
//	// Now from .env
//	ConConf, err := configParser.ParseConnectionConfig(ConConfigPath)
//	if err != nil {
//		log.Fatal("Connection Config Parse error", err)
//	}
//	db, err := database.NewPostgresConnection()
//	if err != nil {
//		log.Fatal("DB connection error", err)
//	}
//
//	FilmRepository := postgres.NewRepo(db)
//	FilmService := services.NewFilmServ(FilmRepository)
//	FilmController := handlers.NewFilmHandler(FilmService)
//	router := restapi.NewRouter(FilmController)
//
//	srv := &http.Server{
//		Addr:    fmt.Sprintf(":" + strconv.Itoa(ConConf.Port)),
//		Handler: router.InitRouter(),
//	}
//	log.Printf("Server is running on port %v\n", ConConf.Port)
//
//	if err := srv.ListenAndServe(); err != nil {
//		log.Fatal("Server error", err)
//	}
//}
