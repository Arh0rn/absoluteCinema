package cache

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
)

const (
	host = "REDIS_HOST"
	port = "REDIS_PORT"
	pass = "REDIS_PASSWORD"
	db   = "REDIS_DB"
)

func NewRedisConnection() (*redis.Client, error) {
	dbIndex, err := strconv.Atoi(os.Getenv(db))
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv(host), os.Getenv(port)),
		Password: os.Getenv(pass),
		DB:       dbIndex,
	})

	return client, nil
}
