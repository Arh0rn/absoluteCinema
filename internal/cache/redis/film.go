package redisCache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Arh0rn/absoluteCinema/pkg/models"
	"github.com/redis/go-redis/v9"
	"time"
)

type FilmCache struct {
	client *redis.Client
}

func (f FilmCache) GetAll(ctx context.Context) ([]*models.Film, error) {
	keys, err := f.client.Keys(ctx, "film*").Result()
	if err != nil {
		return nil, err
	}

	films := make([]*models.Film, 0, len(keys))
	for _, key := range keys {
		data, err := f.client.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}

		var film models.Film
		if err := json.Unmarshal([]byte(data), &film); err != nil {
			return nil, err
		}
		films = append(films, &film)
	}
	return films, nil
}

func (f FilmCache) GetByID(ctx context.Context, id string) (*models.Film, error) {
	key := "film:" + id
	data, err := f.client.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var film models.Film
	if err := json.Unmarshal(data, &film); err != nil {
		return nil, err
	}

	return &film, nil
}

func (f FilmCache) Set(ctx context.Context, film *models.Film) error {
	key := "film:" + film.ID
	data, err := json.Marshal(film)
	if err != nil {
		return err
	}
	cacheTTL, ok := ctx.Value("cacheTTL").(int)
	if !ok {
		return errors.New("cache ttl is not given correctly")
	}

	ttl := time.Duration(cacheTTL) * time.Minute
	return f.client.Set(ctx, key, data, ttl).Err()
}

func (f FilmCache) SetAll(ctx context.Context, films []*models.Film) error {
	for _, film := range films {
		if err := f.Set(ctx, film); err != nil {
			return err
		}
	}
	return nil
}

func (f FilmCache) Update(ctx context.Context, film *models.Film) error {
	err := f.Delete(ctx, film.ID)
	if err != nil {
		return err
	}

	err = f.Set(ctx, film)
	if err != nil {
		return err
	}

	return nil
}

func (f FilmCache) Delete(ctx context.Context, id string) error {
	key := "film:" + id
	return f.client.Del(ctx, key).Err()
}

func NewFilmCache(client *redis.Client) *FilmCache {
	return &FilmCache{
		client: client,
	}
}
