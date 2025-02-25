package pkg

import (
	"github.com/joho/godotenv"
)

//if err := godotenv.Load(); err != nil {
//logrus.Fatalf("error loading env variables: %s", err.Error())
//}
//
//db, err := repository.NewPostgres(repository.Config{
//Host:     viper.GetString("db.host"),
//Port:     viper.GetInt("db.port"),
//User:     viper.GetString("db.user"),
//Password: os.Getenv("POSTGRES_PASSWORD"),
//DBName:   viper.GetString("db.dbname"),
//SSLMode:  viper.GetString("db.sslmode"),
//})

func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}
	return nil
}
