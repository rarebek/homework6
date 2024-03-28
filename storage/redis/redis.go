package redis

import (
	"github.com/gomodule/redigo/redis"
)

type redisRepo struct {
	reds *redis.Pool
}

func NewRedisRepo(reds *redis.Pool) *redisRepo {
	return &redisRepo{reds: reds}
}

func (r *redisRepo) Set(key, value string) error {
	conn := r.reds.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, value)
	return err
}

func (r *redisRepo) Get(key string) (interface{}, error) {
	conn := r.reds.Get()
	defer conn.Close()

	value, err := conn.Do("GET", key)
	return value, err
}

func (r *redisRepo) SetWithTTL(key, value string, seconds int) error {
	conn := r.reds.Get()
	defer conn.Close()
	_, err := conn.Do("SETEX", key, seconds, value)
	return err
}
