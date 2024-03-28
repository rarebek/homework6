package kv

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedisClient(client *redis.Client) *Redis {
	return &Redis{client: client}
}

func (r *Redis) Set(key string, value string, seconds int) error {
	if err := r.client.Set(context.Background(), key, value, time.Duration(seconds)); err != nil {
		return err.Err()
	}

	return nil
}

func (r *Redis) Get(key string) (string, error) {
	str, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}

	return str, nil
}

func (r *Redis) Delete(key string) error {
	if err := r.client.Del(context.Background(), key); err != nil {
		return err.Err()
	}

	return nil
}

func (r *Redis) List() (map[string]string, error) {
	pairs := make(map[string]string)

	cursor := uint64(0)

	for {
		var keys []string
		var err error
		keys, cursor, err := r.client.Scan(context.Background(), cursor, "*", 0).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			value, err := r.client.Get(context.Background(), key).Result()
			if err != nil {
				return nil, err
			}

			pairs[key] = value
		}

		if cursor == 0 {
			break
		}
	}

	return pairs, nil
}
